package processor

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type FakeStore struct {}

func (s *FakeStore) Init() error {
	return nil
}

func (s *FakeStore) Add(filename string, p *pkg.PrintFile) (*pkg.PrintFileRequest, error) {
	return &pkg.PrintFileRequest{
		PrintFile: p,
		Filename:  filename,
		Created:   time.Time{},
		Status:    pkg.Status{},
	}, nil
}

func (s *FakeStore) Update(pfr *pkg.PrintFileRequest) error {
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
	config.SetDefaults()

	assert := assert.New(t)

	printFile := &pkg.PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}
	pj, _ := json.Marshal(printFile)

	s := string(pj)
	fmt.Println(s)

	processor := &Printer{
		&FakeStore{},
		&FakeUpload{},
		&FakeUpload{},
	}
	err := processor.process("test.csv", printFile)
	assert.Nil(err)
}

func createPrintFileEntries(count int) []*pkg.PrintFileEntry {
	entries := make([]*pkg.PrintFileEntry, count)
	for i := 0; i < count; i++ {
		entry := &pkg.PrintFileEntry{
			SampleUnitRef:    "10001",
			Iac:              "ai9bt497r7bn",
			CaseGroupStatus:  "NOTSTARTED",
			EnrolmentStatus:  "",
			RespondentStatus: "",
			Contact: pkg.Contact{
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
