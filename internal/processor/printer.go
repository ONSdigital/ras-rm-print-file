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

func (p *SDCPrinter) Process(filename string, pf *pkg.PrintFile) error {
	log.WithField("filename", filename).Info("processing print file")

	// first save the request to the DB
	err := p.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return err
	}
	printFileRequest, err := p.store.Add(filename, pf)
	if err != nil {
		log.WithError(err).Error("unable to store print file request ")
		return err
	}

	// first sanitise the data
	sanitise(pf)

	// load the ApplyTemplate
	buf, err := applyTemplate(pf)
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

func (p *SDCPrinter) ReProcess(pfr *pkg.PrintFileRequest) error {
	filename := pfr.Filename
	log.WithField("filename", filename).Info("processing print file")

	// first save the request to the DB
	err := p.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return err
	}

	//increment the number of attempts
	numberOfAttempts := pfr.Attempts
	pfr.Attempts = numberOfAttempts + 1

	// load the ApplyTemplate
	buf, err := applyTemplate(pfr.PrintFile)
	if err != nil {
		return err
	}
	pfr.Status.Templated = true
	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("templating complete")

	// first upload to GCS
	if !pfr.Status.UploadedGCS {
		log.Info("print request not uploaded to GCS retrying")
		pfr.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")
	}
	// and then to SFTP
	if !pfr.Status.UploadedSFTP {
		log.Info("print request not uploaded to SFTP retrying")
		pfr.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")
	}

	//check if it's completed
	pfr.Status.Completed = isComplete(pfr)

	err = p.store.Update(pfr)
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
	defer uploader.Close()
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
