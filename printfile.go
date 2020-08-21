package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

//${(actionRequest.sampleUnitRef?trim)!}
//${(actionRequest.iac?trim)!"null"}
//${(actionRequest.caseGroupStatus)!"null"}
//${(actionRequest.enrolmentStatus)!"null"}
//${(actionRequest.respondentStatus)!"null"}
//${(actionRequest.contact.forename?trim)!"null"}:
//${(actionRequest.contact.surname?trim)!"null"}:
//${(actionRequest.contact.emailAddress)!"null"}:
//${(actionRequest.region)!"null"}

var (
	pt = "printfile.tmpl"
)

type PrintFile struct {
	PrintFiles []*PrintFileEntry
}

type PrintFileEntry struct {
	SampleUnitRef string //trim
	Iac string //trim
	CaseGroupStatus string
	EnrolmentStatus string
	RespondentStatus string
	Contact Contact
	Region string

}

type Contact struct {
	Forename string //trim
	Surname string //trim
	EmailAddress string
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

func (pf *PrintFile) process() error {

	//first sanitise the data
	pf.sanitise()

	dat, err := ioutil.ReadFile(pt)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))

	t := template.Must(template.New("printfile.tmpl").ParseFiles(pt))
	//if err != nil {
	//	panic(err)
	//}
	err = t.Execute(os.Stdout, pf)
	if err != nil {
		panic(err)
	}
	return nil
}