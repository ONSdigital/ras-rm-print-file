package sftp

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSFTPAddress(t *testing.T) {
	config.SetDefaults()
	address := createSFTPAddress()
	assert.Equal(t, "localhost:22", address)
}

func TestInitErrorIfNoSFTPServerAvailable(t *testing.T) {
	sftp := SFTPUpload{}
	err := sftp.Init()
	assert.NotNil(t, err)
}

func TestCloseErrorsWithNoConnection(t *testing.T) {
	sftp := SFTPUpload{}
	err := sftp.Close()
	assert.NotNil(t, err)
}

func TestUploadFileErrorsWithNoConnection(t *testing.T) {
	sftp := SFTPUpload{}
	err := sftp.UploadFile("test", []byte("test"))
	assert.NotNil(t, err)
}
