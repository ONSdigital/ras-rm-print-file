package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {
	SetDefaults()
	assert := assert.New(t)

	assert.Equal("debug", viper.GetString("LOG_LEVEL"))
	assert.Equal("ras-rm-printfile", viper.GetString("BUCKET_NAME"))
	assert.Equal("ras-rm-sandbox", viper.GetString("GOOGLE_CLOUD_PROJECT"))
	assert.Equal("localhost", viper.GetString("SFTP_HOST"))
	assert.Equal("22", viper.GetString("SFTP_PORT"))
	assert.Equal("sftp", viper.GetString("SFTP_USERNAME"))
	assert.Equal("sftp", viper.GetString("SFTP_PASSWORD"))
	assert.Equal("printfiles", viper.GetString("SFTP_DIRECTORY"))
}
