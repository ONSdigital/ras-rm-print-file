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

	processor := &SDCPrinter{
		store,
		gcsUpload,
		sftpUpload,
	}
	err := processor.Process("test.csv", printFile)
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

func TestSanitise(t *testing.T) {
	assert := assert.New(t)

	entry := &pkg.PrintFileEntry{
		SampleUnitRef:    "10001 ",
		Iac:              "ai9bt497r7bn ",
		CaseGroupStatus:  "",
		EnrolmentStatus:  "",
		RespondentStatus: "",
		Contact: pkg.Contact{
			Forename:     "",
			Surname:      "",
			EmailAddress: "",
		},
		Region: "",
	}
	entries := []*pkg.PrintFileEntry{entry}
	printFile := &pkg.PrintFile{
		PrintFiles: entries,
	}
	sanitise(printFile)

	assert.Equal("10001", printFile.PrintFiles[0].SampleUnitRef)
	assert.Equal("ai9bt497r7bn", printFile.PrintFiles[0].Iac)
	assert.Equal("null", printFile.PrintFiles[0].CaseGroupStatus)
	assert.Equal("null", printFile.PrintFiles[0].EnrolmentStatus)
	assert.Equal("null", printFile.PrintFiles[0].RespondentStatus)
	assert.Equal("null", printFile.PrintFiles[0].Contact.Forename)
	assert.Equal("null", printFile.PrintFiles[0].Contact.Surname)
	assert.Equal("null", printFile.PrintFiles[0].Contact.EmailAddress)
	assert.Equal("null", printFile.PrintFiles[0].Region)
}

func TestNullIfEmpty(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("null", nullIfEmpty(""))
	assert.Equal("test", nullIfEmpty("test"))
}

func TestApplyTemplate(t *testing.T) {
	assert := assert.New(t)

	printFile := &pkg.PrintFile{
		PrintFiles: createPrintFileEntries(1),
	}

	buffer, err := applyTemplate(printFile)
	assert.Nil(err)
	assert.Equal("10001:ai9bt497r7bn:NOTSTARTED:::Jon:Snow:jon.snow@example.com:HH\n", buffer.String())
}

func TestApplyTemplateEmptyPrintFile(t *testing.T) {
	assert := assert.New(t)

	printFile := &pkg.PrintFile{}

	buffer, err := applyTemplate(printFile)
	assert.Nil(err)
	assert.Equal("", buffer.String())
}
