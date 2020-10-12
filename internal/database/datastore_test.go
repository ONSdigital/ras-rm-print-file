package database

import (
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddErrorsWithoutInit(t *testing.T) {
	ds := DataStore{}
	assert := assert.New(t)

	// if there's no client initialise then this method should error immediately
	pfr, err := ds.Add("test.csv", "test.json")
	assert.NotNil(err)
	assert.Nil(pfr)
}

func TestUpload(t *testing.T) {
	ds := DataStore{}
	// if there's no client initialise then this method should error immediately
	err := ds.Update(&pkg.PrintFileRequest{})
	assert.NotNil(t, err)
}

func TestFindIncomplete(t *testing.T) {
	assert := assert.New(t)
	ds := DataStore{}
	// if there's no client initialise then this method should error immediately
	pfr, err := ds.FindIncomplete()
	assert.NotNil(err)
	assert.Nil(pfr)
}

func TestCreatePrintFileRequest(t *testing.T) {
	assert := assert.New(t)
	printfile := "test.csv"
	datafile := "test.json"
	pfr := createPrintFileRequest(printfile, datafile)
	assert.Equal(printfile, pfr.PrintFilename)
	assert.Equal(datafile, pfr.DataFilename)
	assert.Equal(1, pfr.Attempts)
	assert.False(pfr.Status.UploadedSFTP)
	assert.False(pfr.Status.UploadedGCS)
	assert.False(pfr.Status.Templated)
	assert.False(pfr.Status.Completed)
	assert.NotNil(pfr.Created)
	assert.NotNil(pfr.Updated)
}
