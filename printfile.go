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
	SampleUnitRef string `json:"sampleUnitRef"`
	Iac string `json:"iac"`
	CaseGroupStatus string `json:"caseGroupStatus"`
	EnrolmentStatus string `json:"enrolmentStatus"`
	RespondentStatus string `json:"respondentStatus"`
	Contact Contact `json:"contact"`
	Region string `json:"region"`

}

type Contact struct {
	Forename string `json:"forename"`
	Surname string `json:"surname"`
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

func (pf *PrintFile) process(filename string) error {
	log.WithField("filename", filename).Info("processing print file")
	// first save the request to the DB
	store := &Store{}
	store.store(filename, pf)

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
	log.WithField("template", printTemplate).WithField("filename", filename).Info("templating complete")
	pf.upload(filename, buf)

	//TODO handle errors/retry
	return nil
}

func (pf *PrintFile) upload(filename string, buffer *bytes.Buffer) {
	log.WithField("filename", filename).Info("uploading file")
	// first upload to GCS
	gcsUpload := &GCSUpload{}
	gcsUpload.Init()
	err := gcsUpload.UploadFile(filename, buffer.Bytes())
	if err != nil {
		//TODO retry
	}

	// and then to SFTP
	sftpUpload := SFTPUpload{}
	err = sftpUpload.Init()
	if err != nil {
		// TODO retry connection ?
		return
	}
	err = sftpUpload.UploadFile(filename, buffer.Bytes())
	if err != nil {
		//TODO retry
	}
}