package processor

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	mocks "github.com/ONSdigital/ras-rm-print-file/mocks/pkg"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestProcess(t *testing.T) {
	config.SetDefaults()

	assert := assert.New(t)

	printFile := &pkg.PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("Add", mock.Anything, mock.Anything).Return(&pkg.PrintFileRequest{}, nil)
	store.On("Update", mock.Anything).Return(nil)

	gcsUpload := new(mocks.Upload)
	gcsUpload.On("Init").Return(nil)
	gcsUpload.On("UploadFile", mock.Anything, mock.Anything).Return(nil)

	sftpUpload := new(mocks.Upload)
	sftpUpload.On("Init").Return(nil)
	sftpUpload.On("UploadFile", mock.Anything, mock.Anything).Return(nil)

	processor := &Printer{
		store,
		gcsUpload,
		sftpUpload,
	}
	err := processor.process("test.csv", printFile)
	assert.Nil(err)

	store.AssertExpectations(t)
	gcsUpload.AssertExpectations(t)
	sftpUpload.AssertExpectations(t)
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
