package processor

import (
	"bytes"
	"fmt"
	"github.com/ONSdigital/ras-rm-print-file/internal/database"
	"github.com/ONSdigital/ras-rm-print-file/internal/gcs"
	"github.com/ONSdigital/ras-rm-print-file/internal/sftp"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

var (
	printTemplate = "printfile.tmpl"
)

type Printer struct {
	store      pkg.Store
	gcsUpload  pkg.Upload
	sftpUpload pkg.Upload
}


func Process(filename string, printFile *pkg.PrintFile) error {
	processor := &Printer{}
	processor.store = &database.DataStore{}
	processor.gcsUpload = &gcs.GCSUpload{}
	processor.sftpUpload = &sftp.SFTPUpload{}
	return processor.process(filename, printFile)
}

func (p *Printer) process(filename string, printFile *pkg.PrintFile) error {
	log.WithField("filename", filename).Info("processing print file")

	// first save the request to the DB
	err := p.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return err
	}
	pfr, err := p.store.Add(filename, printFile)
	if err != nil {
		log.WithError(err).Error("unable to store print file request ")
		return err
	}

	// first sanitise the data
	sanitise(printFile)

	// load the ApplyTemplate
	buf, err := applyTemplate(printFile, filename)
	if err != nil {
		return err
	}
	pfr.Status.TemplateComplete = true
	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("templating complete")

	// first upload to GCS
	pfr.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")

	// and then to SFTP
	pfr.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")

	err = p.store.Update(pfr)
	if err != nil {
		log.WithError(err).Error("failed to Update database")
		//TODO set to not ready
		return err
	}
	return nil
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
		fmt.Print(pfe)
	}
}

func nullIfEmpty(value string) string {
	if value == "" {
		log.WithField("value", value).Debug("empty value replacing with null")
		return "null"
	}
	return value
}

func applyTemplate(pf *pkg.PrintFile, filename string) (*bytes.Buffer, error) {
	// load the ApplyTemplate
	log.WithField("ApplyTemplate", printTemplate).Info("about to load ApplyTemplate")
	t, err := template.New(printTemplate).ParseFiles(printTemplate)
	if err != nil {
		log.WithError(err).Error("failed to find ApplyTemplate")
		//TODO set to not ready
		return nil, err
	}

	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("about to process ApplyTemplate")
	// create a bytes buffer and run the ApplyTemplate engine
	buf := &bytes.Buffer{}
	err = t.Execute(buf, pf)
	if err != nil {
		log.WithError(err).Error("failed to process ApplyTemplate")
		return nil, err
	}
	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("templating complete")
	return buf, nil
}

