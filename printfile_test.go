package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintFile(t *testing.T) {
	setDefaults()

	assert := assert.New(t)

	printfile := &PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}
	pj, _ := json.Marshal(printfile)

	s := string(pj)
	fmt.Println(s)
	err := printfile.process("test.csv")
	assert.Nil(err)
}

func createPrintFileEntries(count int) []*PrintFileEntry {
	entries := make([]*PrintFileEntry, count)
	for i := 0; i < count; i++ {
		entry := &PrintFileEntry{
			SampleUnitRef:    "10001",
			Iac:              "ai9bt497r7bn",
			CaseGroupStatus:  "NOTSTARTED",
			EnrolmentStatus:  "",
			RespondentStatus: "",
			Contact: Contact{
				Forename:     "Jon",
				Surname:      "Snow",
				EmailAddress: "jon.snow@example.com",
			},
			Region: "HH",
		}
		entries[i] = entry
	}
	return entries
}
