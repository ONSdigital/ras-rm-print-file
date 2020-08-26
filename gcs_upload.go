package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type GCSUpload struct {
	ctx context.Context
	client *storage.Client
}

func (u *GCSUpload) Init() error {
	var err error
	u.ctx = context.Background()
	u.client, err = storage.NewClient(u.ctx)
	if err != nil {
		log.WithError(err).Error("error creating gcp client")
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	return nil
}

func (u *GCSUpload) Close() {
	u.client.Close()
}

// uploadFile uploads an object.
func (u *GCSUpload) UploadFile(name string, contents []byte) error {
	bucket := viper.GetString("BUCKET_NAME")

	ctx, cancel := context.WithTimeout(u.ctx, time.Second*50)
	defer cancel()

	// GCSUpload an object with storage.Writer.
	wc := u.client.Bucket(bucket).Object(name).NewWriter(ctx)
	if _, err := wc.Write(contents); err != nil {
		log.WithError(err).Error("error writing bytes to bucket")
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.WithError(err).Error("error closing bucket writer")
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}