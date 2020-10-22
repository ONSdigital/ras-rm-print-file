package gcs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GCSUpload struct {
	ctx    context.Context
	client *storage.Client
}

func (u *GCSUpload) Init() error {
	var err error
	u.ctx = context.Background()
	log.Info("initialising GCS connection")
	u.client, err = storage.NewClient(u.ctx)
	if err != nil {
		log.WithError(err).Error("error creating gcp client")
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	log.Info("connected to GCS")
	return nil
}

func (u *GCSUpload) Close() error {
	if u.client == nil {
		return errors.New("please initialise the connection")
	}
	log.Info("closing connection to GCS")
	err := u.client.Close()
	log.Info("GCS connection closed")
	return err
}

// uploadFile uploads an object.
func (u *GCSUpload) UploadFile(filename string, contents []byte) error {
	if u.client == nil {
		return errors.New("please initialise the connection")
	}
	bucket := viper.GetString("BUCKET_NAME")
	path := bucketPath(filename)

	log.WithField("filename", path).WithField("bucket", bucket).Info("uploading to bucket")

	ctx, cancel := context.WithTimeout(u.ctx, time.Second*50)
	defer cancel()

	// GCSUpload an object with storage.Writer.
	wc := u.client.Bucket(bucket).Object(path).NewWriter(ctx)
	log.WithField("filename", path).WithField("bucket", bucket).Info("about to write contents to bucket")
	if _, err := wc.Write(contents); err != nil {
		log.WithError(err).WithField("bucket", bucket).WithField("path", path).Error("error writing bytes to bucket ")
		return err
	}
	if err := wc.Close(); err != nil {
		log.WithError(err).Error("error closing bucket writer")
		return err
	}
	log.WithField("filename", path).WithField("bucket", bucket).Info("upload to bucket complete")
	return nil
}
