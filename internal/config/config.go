package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConfigureLogging() {
	logLevel := viper.GetString("LOG_LEVEL")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithError(err).WithField("logLevel", logLevel).Error("invalid log level")
		//default to debug
		level = log.DebugLevel
	}
	log.SetLevel(level)
	log.WithField("level", logLevel).Debug("log level set")
}

func SetDefaults() {
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("BUCKET_NAME", "ras-rm-printfile")
	viper.SetDefault("GOOGLE_CLOUD_PROJECT", "ras-rm-sandbox")
	viper.SetDefault("SFTP_HOST", "localhost")
	viper.SetDefault("SFTP_PORT", "22")
	viper.SetDefault("SFTP_USERNAME", "sftp")
	viper.SetDefault("SFTP_PASSWORD", "sftp")
	viper.SetDefault("SFTP_DIRECTORY", "printfiles")
	viper.SetDefault("RETRY_DELAY", "60000")
}
