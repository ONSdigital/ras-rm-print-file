package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type FakeStore struct {}

func (s *FakeStore) Init() error {
	return nil
}

func (s *FakeStore) store(filename string, p *PrintFile) (*PrintFileRequest, error) {
	return &PrintFileRequest{
		printFile: p,
		filename:  filename,
		created:   time.Time{},
		Status:    Status{},
	}, nil
}

func (s *FakeStore) update(pfr *PrintFileRequest) error {
	return nil
}

func TestPrintFile(t *testing.T) {
	setDefaults()

	assert := assert.New(t)

	printFile := &PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}
	pj, _ := json.Marshal(printFile)

	s := string(pj)
	fmt.Println(s)
	err := printFile.process(&FakeStore{}, "test.csv")
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
