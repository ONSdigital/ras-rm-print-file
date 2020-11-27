package retry

import (
	mocks "github.com/ONSdigital/ras-rm-print-file/mocks/pkg"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {

	printFileRequest := &pkg.PrintFileRequest{
		DataFilename: "test.json",
		PrintFilename:  "test.csv",
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
	store.On("Close").Return(nil)

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

	printFileRequest := &pkg.PrintFileRequest{
		DataFilename: "test.json",
		PrintFilename:  "test.csv",
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
	store.On("Close").Return(nil)

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
