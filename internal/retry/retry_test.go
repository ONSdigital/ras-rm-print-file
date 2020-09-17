package retry

import (
	mocks "github.com/ONSdigital/ras-rm-print-file/mocks/pkg"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	printFileEntry := &pkg.PrintFileEntry{
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

	printFile := &pkg.PrintFile{
		PrintFiles: []*pkg.PrintFileEntry{printFileEntry},
	}

	printFileRequest := &pkg.PrintFileRequest{
		PrintFile: printFile,
		Filename:  "test.csv",
		Created:   time.Now(),
		Status: pkg.Status{
			Templated:    true,
			UploadedGCS:  false,
			UploadedSFTP: false,
		},
	}

	printFileRequests := []*pkg.PrintFileRequest{printFileRequest}

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("FindIncomplete", mock.Anything, mock.Anything).Return(printFileRequests, nil)

	printer := new(mocks.Printer)
	printer.On("ReProcess", printFileRequest).Return(nil)

	backOffRetry := BackoffRetry{
		store:   store,
		printer: printer,
	}
	backOffRetry.process()

	store.AssertNotCalled(t, "Update")
	store.AssertNotCalled(t, "Add")
	store.AssertExpectations(t)

	printer.AssertCalled(t, "ReProcess", printFileRequest)
}

func TestReProcessWhenCompleteDoesNotRun(t *testing.T) {
	printFileEntry := &pkg.PrintFileEntry{
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

	printFile := &pkg.PrintFile{
		PrintFiles: []*pkg.PrintFileEntry{printFileEntry},
	}

	printFileRequest := &pkg.PrintFileRequest{
		PrintFile: printFile,
		Filename:  "test.csv",
		Created:   time.Now(),
		Status: pkg.Status{
			Templated:    true,
			UploadedGCS:  true,
			UploadedSFTP: true,
		},
	}

	printFileRequests := []*pkg.PrintFileRequest{printFileRequest}

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("FindIncomplete", mock.Anything, mock.Anything).Return(printFileRequests, nil)

	printer := new(mocks.Printer)
	printer.On("ReProcess", printFileRequest).Return(nil)

	backOffRetry := BackoffRetry{
		store:   store,
		printer: printer,
	}
	backOffRetry.process()

	// need to wait a few milliseconds for the go rountine to execture
	time.Sleep(10 * time.Millisecond)

	store.AssertNotCalled(t, "Update")
	store.AssertNotCalled(t, "Add")
	store.AssertExpectations(t)

	printer.AssertCalled(t, "ReProcess", printFileRequest)
}
