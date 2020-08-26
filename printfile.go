package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"text/template"
)

var (
	pt = "printfile.tmpl"
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
	fmt.Printf("before %q\n", value)
	if value == "" {
		fmt.Print("return null\n")
		return "null"
	}
	return value
}

func (pf *PrintFile) process(filename string) error {
	//first save the request to the DB
	store := &Store{}
	store.store(filename, pf)

	//first sanitise the data
	pf.sanitise()

	dat, err := ioutil.ReadFile(pt)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))

	t, err := template.New("printfile.tmpl").ParseFiles(pt)
	if err != nil {
		log.WithError(err).Error("failed to find template")
		return err
	}
	fmt.Println(pf)

	buf := &bytes.Buffer{}
	err = t.Execute(buf, pf)
	if err != nil {
		log.WithError(err).Error("failed to process template")
	}
	fmt.Println(buf.String())
	upload(filename, buf)

	//TODO handle errors/retry
	return nil
}

func upload(filename string, buffer *bytes.Buffer) {
	// first upload to GCS
	gcsUpload := &GCSUpload{}
	gcsUpload.Init()
	gcsUpload.UploadFile(filename, buffer.Bytes())

	// and then to SFTP
	sftpUpload := SFTPUpload{}
	sftpUpload.Init()
	sftpUpload.UploadFile(filename, buffer.Bytes())
}