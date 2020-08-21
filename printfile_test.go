package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintFile(t *testing.T) {

	assert := assert.New(t)

	printfile := &PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}
	fmt.Println(printfile)
	err := printfile.process()
	assert.Nil(err)
}

func createPrintFileEntries(count int) []*PrintFileEntry {
	entries := make([]*PrintFileEntry, count)
	for i := 0; i < count; i ++ {
		entry := &PrintFileEntry{
			SampleUnitRef:    "10001",
			Iac:              "ai9bt497r7bn",
			CaseGroupStatus:  "NOTSTARTED",
			EnrolmentStatus:  "",
			RespondentStatus: "",
			Contact:          Contact{
				Forename:     "Jon",
				Surname:      "Snow",
				EmailAddress: "jon.snow@example.com",
			},
			Region:           "HH",
		}
		entries[i] = entry
	}
	return entries
}