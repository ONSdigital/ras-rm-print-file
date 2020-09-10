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

func (s *FakeStore) Add(filename string, p *PrintFile) (*PrintFileRequest, error) {
	return &PrintFileRequest{
		printFile: p,
		filename:  filename,
		created:   time.Time{},
		Status:    Status{},
	}, nil
}

func (s *FakeStore) Update(pfr *PrintFileRequest) error {
	return nil
}

type FakeUpload struct {}

func (u *FakeUpload) Init() error {
	return nil
}

func (u *FakeUpload) Close() {

}
func (u *FakeUpload) UploadFile(filename string, contents []byte) error {
	return nil
}

func TestProcess(t *testing.T) {
	setDefaults()

	assert := assert.New(t)

	printFile := &PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}
	pj, _ := json.Marshal(printFile)

	s := string(pj)
	fmt.Println(s)

	processor := &Processor{
		&FakeStore{},
		&FakeUpload{},
		&FakeUpload{},
	}
	err := processor.process("test.csv", printFile)
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
