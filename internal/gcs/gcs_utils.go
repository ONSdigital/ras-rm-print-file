package gcs

import (
	"os"

	"github.com/spf13/viper"
)

func bucketPath(filename string) string {
	prefix := viper.GetString("PREFIX_NAME")

	path := filename
	if prefix != "" {
		ps := string(os.PathSeparator)
		path = ps + prefix + ps + filename
	}

	return path
}
