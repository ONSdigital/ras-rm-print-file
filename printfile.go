package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

var (
	printTemplate = "printfile.tmpl"
)

type PrintFile struct {
	PrintFiles []*PrintFileEntry
}


type PrintFileEntry struct {
	SampleUnitRef    string  `json:"sampleUnitRef"`
	Iac              string  `json:"iac"`
	CaseGroupStatus  string  `json:"caseGroupStatus"`
	EnrolmentStatus  string  `json:"enrolmentStatus"`
	RespondentStatus string  `json:"respondentStatus"`
	Contact          Contact `json:"contact"`
	Region           string  `json:"region"`
}

type Contact struct {
	Forename     string `json:"forename"`
	Surname      string `json:"surname"`
	EmailAddress string `json:"emailAddress"`
}

func (pf *PrintFile) sanitise() {
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

func (pf *PrintFile) process(str Store, filename string) error {
	log.WithField("filename", filename).Info("processing print file")
	// first save the request to the DB
	str.Init()
	pfr, err := str.store(filename, pf)
	if err != nil {
		log.WithError(err).Error("unable to store print file request ")
		return err
	}

	// first sanitise the data
	pf.sanitise()

	// load the template
	log.WithField("template", printTemplate).Info("about to load template")
	t, err := template.New(printTemplate).ParseFiles(printTemplate)
	if err != nil {
		log.WithError(err).Error("failed to find template")
		//TODO set to not ready
		return err
	}

	log.WithField("template", printTemplate).WithField("filename", filename).Info("about to process template")
	// create a bytes buffer and run the template engine
	buf := &bytes.Buffer{}
	err = t.Execute(buf, pf)
	if err != nil {
		log.WithError(err).Error("failed to process template")
		return nil
	}
	pfr.Status.TemplateComplete = true
	log.WithField("template", printTemplate).WithField("filename", filename).Info("templating complete")

	err = pf.uploadGCS(filename, buf)
	if err != nil {
		pfr.Status.UploadedGCS = false
	} else {
		pfr.Status.UploadedGCS = true
	}
	err = pf.uploadSFTP(filename, buf)
	if err != nil {
		pfr.Status.UploadedGCS = false
	} else {
		pfr.Status.UploadedGCS = true
	}

	err = str.update(pfr)
	if err != nil {
		log.WithError(err).Error("failed to update database")
		//TODO set to not ready
		return err
	}
	return nil
}

func (pf *PrintFile) uploadGCS(filename string, buffer *bytes.Buffer) error {
	log.WithField("filename", filename).Info("uploading file to gcs")
	// first upload to GCS
	gcsUpload := &GCSUpload{}
	err := gcsUpload.Init()
	if err != nil {
		return err
	}
	return gcsUpload.UploadFile(filename, buffer.Bytes())
}

func (pf *PrintFile) uploadSFTP(filename string, buffer *bytes.Buffer) error {
	log.WithField("filename", filename).Info("uploading file to sftp")
	// and then to SFTP
	sftpUpload := SFTPUpload{}
	err := sftpUpload.Init()
	if err != nil {
		return err
	}
	return sftpUpload.UploadFile(filename, buffer.Bytes())
}
