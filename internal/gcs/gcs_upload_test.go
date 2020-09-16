package gcs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloseErrorsWithNoConnection(t *testing.T) {
	gcs := GCSUpload{}
	err := gcs.Close()
	assert.NotNil(t, err)
}

func TestUploadFileErrorsWithNoConnection(t *testing.T) {
	gcs := GCSUpload{}
	err := gcs.UploadFile("test", []byte("test"))
	assert.NotNil(t, err)
}