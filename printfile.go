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

func (pf *PrintFile) ApplyTemplate(filename string) (*bytes.Buffer, error) {
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
