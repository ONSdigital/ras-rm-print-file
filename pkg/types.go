package pkg

import "time"

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

type PrintFileRequest struct {
	PrintFile *PrintFile
	Filename  string
	Created   time.Time
	Status    Status
}

type Status struct {
	TemplateComplete bool
	UploadedGCS      bool
	UploadedSFTP     bool
}
