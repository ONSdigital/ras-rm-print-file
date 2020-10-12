package gcs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownloadCloseErrorsWithNoConnection(t *testing.T) {
	gcs := GCSDownload{}
	err := gcs.Close()
	assert.NotNil(t, err)
}

func TestDownloadUploadFileErrorsWithNoConnection(t *testing.T) {
	gcs := GCSDownload{}
	printfile, err := gcs.DownloadFile("test")
	assert.NotNil(t, err)
	assert.Nil(t, printfile)
}

