package gcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("someFileName.json", bucket_path("someFileName.json"))
}
