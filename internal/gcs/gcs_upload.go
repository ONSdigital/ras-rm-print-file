package gcs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GCSUpload struct {
	ctx    context.Context
	client *storage.Client
}

func (u *GCSUpload) Init() error {
	var err error
	u.ctx = context.Background()
	logger.Info("initialising GCS connection")
	u.client, err = storage.NewClient(u.ctx)
	if err != nil {
		logger.Error("error creating gcp client",
			zap.Error(err))
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	logger.Info("connected to GCS")
	return nil
}

func (u *GCSUpload) Close() error {
	if u.client == nil {
		return errors.New("please initialise the connection")
	}
	logger.Info("closing connection to GCS")
	err := u.client.Close()
	logger.Info("GCS connection closed")
	return err
}

// uploadFile uploads an object.
func (u *GCSUpload) UploadFile(filename string, contents []byte) error {
	if u.client == nil {
		return errors.New("please initialise the connection")
	}
	bucket := viper.GetString("BUCKET_NAME")
	path := bucketPath(filename)

	logger.Info("uploading to bucket",
		zap.String("filename", path),
		zap.String("bucket", bucket))

	ctx, cancel := context.WithTimeout(u.ctx, time.Second*50)
	defer cancel()

	// GCSUpload an object with storage.Writer.
	wc := u.client.Bucket(bucket).Object(path).NewWriter(ctx)
	logger.Info("about to write contents to bucket",
		zap.String("filename", path),
		zap.String("bucket", bucket))

	if _, err := wc.Write(contents); err != nil {
		logger.Error("error writing bytes to bucket",
			zap.String("bucket", bucket),
			zap.String("path", path))
		return err
	}
	if err := wc.Close(); err != nil {
		logger.Error("error closing bucket writer",
			zap.Error(err))
		return err
	}
	logger.Info("upload to bucket complete",
		zap.String("filename", path),
		zap.String("bucket", bucket))
	return nil
}
