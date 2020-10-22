package gcs

import (
	"os"

	"github.com/spf13/viper"
)

func bucket_path(filename string) string {
	prefix := viper.GetString("PREFIX_NAME")

	path := filename
	if prefix != "" {
		ps := string(os.PathSeparator)
		path = prefix + ps + filename
	}

	return path
}
