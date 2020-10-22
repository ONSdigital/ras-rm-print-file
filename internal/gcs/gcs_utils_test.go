package gcs

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("someFileName.json", bucketPath("someFileName.json"))

	viper.SetDefault("PREFIX_NAME", "test")
	assert.Equal("test/someFileName.json", bucketPath("someFileName.json"))

}
