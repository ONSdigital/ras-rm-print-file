package config

import (
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("BUCKET_NAME", "ras-rm-print-file")
	viper.SetDefault("PREFIX_NAME", "")
	viper.SetDefault("GOOGLE_CLOUD_PROJECT", "ras-rm-sandbox")
	viper.SetDefault("SFTP_HOST", "localhost")
	viper.SetDefault("SFTP_PORT", "22")
	viper.SetDefault("SFTP_USERNAME", "sftp")
	viper.SetDefault("SFTP_PASSWORD", "sftp")
	viper.SetDefault("SFTP_DIRECTORY", "printfiles")
	viper.SetDefault("RETRY_DELAY", "3600000")
	viper.SetDefault("CLEANUP_DELAY", "24")
	viper.SetDefault("CLEANUP_RETENTION", "720")
	viper.SetDefault("PUBSUB_SUB_ID", "print-file")
	viper.SetDefault("PUBSUB_TOPIC", "print-file")
}
