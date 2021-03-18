package database

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	mocks "github.com/ONSdigital/ras-rm-print-file/mocks/pkg"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCleanup(t *testing.T) {
	config.SetDefaults()

	createdUpdatedTime := time.Now().AddDate(0, 0, -31)

	printFileRequest := &pkg.PrintFileRequest{
		DataFilename:  "test.json",
		PrintFilename: "test.csv",
		Created:       createdUpdatedTime,
		Updated:       createdUpdatedTime,
		Status: pkg.Status{
			Templated:    true,
			UploadedGCS:  true,
			UploadedSFTP: true,
			Completed:    true,
		},
	}

	printFileRequests := []*pkg.PrintFileRequest{printFileRequest}

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("FindComplete", mock.Anything, mock.Anything).Return(printFileRequests, nil)
	store.On("Close").Return(nil)
	store.On("Delete", printFileRequest).Return(nil)

	cleanUp := CleanUp{
		store: store,
	}
	cleanUp.process()

	store.AssertNotCalled(t, "Update")
	store.AssertNotCalled(t, "Add")
	store.AssertExpectations(t)

}

func TestCleanUpDoesNotRunWhenDurationLessThan30days(t *testing.T) {
	config.SetDefaults()

	printFileRequest := &pkg.PrintFileRequest{
		DataFilename:  "test.json",
		PrintFilename: "test.csv",
		Created:       time.Now(),
		Updated:       time.Now(),
		Status: pkg.Status{
			Templated:    true,
			UploadedGCS:  true,
			UploadedSFTP: true,
			Completed:    true,
		},
	}

	printFileRequests := []*pkg.PrintFileRequest{printFileRequest}

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("FindComplete", mock.Anything, mock.Anything).Return(printFileRequests, nil)
	store.On("Close").Return(nil)

	cleanUp := CleanUp{
		store: store,
	}
	cleanUp.process()

	store.AssertNotCalled(t, "Update")
	store.AssertNotCalled(t, "Add")
	store.AssertNotCalled(t, "Delete")
	store.AssertExpectations(t)
}

func TestCleanUpDoesNotRunWhenNotComplete(t *testing.T) {
	config.SetDefaults()

	printFileRequest := &pkg.PrintFileRequest{
		DataFilename:  "test.json",
		PrintFilename: "test.csv",
		Created:       time.Now(),
		Updated:       time.Now(),
		Status: pkg.Status{
			Templated:    true,
			UploadedGCS:  true,
			UploadedSFTP: false,
			Completed:    false,
		},
	}

	printFileRequests := []*pkg.PrintFileRequest{printFileRequest}

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("FindComplete", mock.Anything, mock.Anything).Return(printFileRequests, nil)
	store.On("Close").Return(nil)

	cleanUp := CleanUp{
		store: store,
	}
	cleanUp.process()

	store.AssertNotCalled(t, "Update")
	store.AssertNotCalled(t, "Add")
	store.AssertNotCalled(t, "Delete")
	store.AssertExpectations(t)
}

func TestCleanUpDoesNotRunWhenNoResults(t *testing.T) {
	config.SetDefaults()

	store := new(mocks.Store)
	store.On("Init").Return(nil)
	store.On("FindComplete", mock.Anything, mock.Anything).Return([]*pkg.PrintFileRequest{}, nil)
	store.On("Close").Return(nil)

	cleanUp := CleanUp{
		store: store,
	}
	cleanUp.process()

	store.AssertNotCalled(t, "Update")
	store.AssertNotCalled(t, "Add")
	store.AssertNotCalled(t, "Delete")
	store.AssertExpectations(t)
}
