package processor

import (
	"bytes"
	"github.com/ONSdigital/ras-rm-print-file/internal/database"
	"github.com/ONSdigital/ras-rm-print-file/internal/gcs"
	"github.com/ONSdigital/ras-rm-print-file/internal/sftp"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"text/template"
)

var printTemplate = "printfile.tmpl"

type SDCPrinter struct {
	store      pkg.Store
	gcsUpload  pkg.Upload
	sftpUpload pkg.Upload
}

func Create() *SDCPrinter {
	processor := &SDCPrinter{}
	processor.store = &database.DataStore{}
	processor.gcsUpload = &gcs.GCSUpload{}
	processor.sftpUpload = &sftp.SFTPUpload{}
	return processor
}

func CreateAndProcess(filename string, printFile *pkg.PrintFile) error {
	processor := Create()
	return processor.Process(filename, printFile)
}

func (p *SDCPrinter) Process(filename string, printFile *pkg.PrintFile) error {
	log.WithField("filename", filename).Info("processing print file")

	// first save the request to the DB
	err := p.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return err
	}
	printFileRequest, err := p.store.Add(filename, printFile)
	if err != nil {
		log.WithError(err).Error("unable to store print file request ")
		return err
	}

	// first sanitise the data
	sanitise(printFile)

	// load the ApplyTemplate
	buf, err := applyTemplate(printFile)
	if err != nil {
		return err
	}
	printFileRequest.Status.Templated = true
	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("templating complete")

	// first upload to GCS
	printFileRequest.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")

	// and then to SFTP
	printFileRequest.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")

	//check if it's completed
	printFileRequest.Status.Completed = isComplete(printFileRequest)

	err = p.store.Update(printFileRequest)
	if err != nil {
		log.WithError(err).Error("failed to Update database")
		return err
	}
	return nil
}

func (p *SDCPrinter) ReProcess(printFileRequest *pkg.PrintFileRequest) error {
	filename := printFileRequest.Filename
	log.WithField("filename", filename).Info("processing print file")

	// first save the request to the DB
	err := p.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return err
	}

	//increment the number of attempts
	numberOfAttempts := printFileRequest.Attempts
	printFileRequest.Attempts = numberOfAttempts + 1

	// load the ApplyTemplate
	buf, err := applyTemplate(printFileRequest.PrintFile)
	if err != nil {
		return err
	}
	printFileRequest.Status.Templated = true
	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("templating complete")

	// first upload to GCS
	if !printFileRequest.Status.UploadedGCS {
		log.Info("print request not uploaded to GCS retrying")
		printFileRequest.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")
	}
	// and then to SFTP
	if !printFileRequest.Status.UploadedSFTP {
		log.Info("print request not uploaded to SFTP retrying")
		printFileRequest.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")
	}

	//check if it's completed
	printFileRequest.Status.Completed = isComplete(printFileRequest)

	err = p.store.Update(printFileRequest)
	if err != nil {
		log.WithError(err).Error("failed to Update database")
		return err
	}
	return nil
}

func isComplete(printFileRequest *pkg.PrintFileRequest) bool {
	return printFileRequest.Status.Templated && printFileRequest.Status.UploadedGCS && printFileRequest.Status.UploadedSFTP
}

func upload(filename string, buffer *bytes.Buffer, uploader pkg.Upload, name string) bool {
	log.WithField("filename", filename).Infof("uploading file to %v", name)
	err := uploader.Init()
	if err != nil {
		log.WithError(err).Errorf("failed to initialise %v upload", name)
		return false
	}
	err = uploader.UploadFile(filename, buffer.Bytes())
	if err != nil {
		log.WithError(err).Errorf("failed to upload to %v", name)
		return false
	}
	return true
}

func sanitise(pf *pkg.PrintFile) {
	log.Info("sanitising print file to match expected outcomes")
	for _, pfe := range pf.PrintFiles {
		pfe.SampleUnitRef = strings.TrimSpace(pfe.SampleUnitRef)
		pfe.Iac = nullIfEmpty(strings.TrimSpace(pfe.Iac))
		pfe.CaseGroupStatus = nullIfEmpty(pfe.CaseGroupStatus)
		pfe.EnrolmentStatus = nullIfEmpty(pfe.EnrolmentStatus)
		pfe.RespondentStatus = nullIfEmpty(pfe.RespondentStatus)
		pfe.Contact.Forename = nullIfEmpty(pfe.Contact.Forename)
		pfe.Contact.Surname = nullIfEmpty(pfe.Contact.Surname)
		pfe.Contact.EmailAddress = nullIfEmpty(pfe.Contact.EmailAddress)
		pfe.Region = nullIfEmpty(pfe.Region)
	}
}

func nullIfEmpty(value string) string {
	if value == "" {
		log.Debug("empty value replacing with null")
		return "null"
	}
	return value
}

func applyTemplate(pf *pkg.PrintFile) (*bytes.Buffer, error) {
	//find the template
	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Error("unable to load template")
	}
	// Template location
	templateLocation := wd + "/templates/" + printTemplate

	// load the ApplyTemplate
	log.WithField("ApplyTemplate", printTemplate).Info("about to load ApplyTemplate")
	t, err := template.New(printTemplate).ParseFiles(templateLocation)
	if err != nil {
		log.WithError(err).Error("failed to find ApplyTemplate")
		return nil, err
	}

	log.WithField("ApplyTemplate", printTemplate).Info("about to process ApplyTemplate")
	// create a bytes buffer and run the ApplyTemplate engine
	buf := &bytes.Buffer{}
	err = t.Execute(buf, pf)
	if err != nil {
		log.WithError(err).Error("failed to process ApplyTemplate")
		return nil, err
	}
	log.WithField("ApplyTemplate", printTemplate).Info("templating complete")
	return buf, nil
}
