package gcs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GCSDownload struct {
	ctx    context.Context
	client *storage.Client
}

func (d *GCSDownload) Init() error {
	var err error
	d.ctx = context.Background()
	log.Info("initialising GCS connection")
	d.client, err = storage.NewClient(d.ctx)
	if err != nil {
		log.WithError(err).Error("error creating gcp client")
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	log.Info("connected to GCS")
	return nil
}

func (d *GCSDownload) Close() error {
	if d.client == nil {
		return errors.New("please initialise the connection")
	}
	log.Info("closing connection to GCS")
	err := d.client.Close()
	log.Info("GCS connection closed")
	return err
}

func (d *GCSDownload) DownloadFile(filename string) (*pkg.PrintFile, error) {
	// loads the payload from a GCS bucket
	if d.client == nil {
		return nil, errors.New("please initialise the connection")
	}
	bucket := viper.GetString("BUCKET_NAME")
	prefix := viper.GetString("PREFIX_NAME")
	path := prefix + filename
	log.WithField("filename", path).WithField("bucket", bucket).Info("downloading from bucket")

	ctx, cancel := context.WithTimeout(d.ctx, time.Second*50)
	defer cancel()

	// GCSUpload an object with storage.Writer.
	rc, err := d.client.Bucket(bucket).Object(path).NewReader(ctx)
	if err != nil {
		log.WithError(err).Error("error reading from bucket " + path)
		return nil, err
	}
	log.WithField("filename", path).WithField("bucket", bucket).Info("about to read contents from bucket")

	buf := &bytes.Buffer{}
	defer rc.Close()
	if _, err := io.Copy(buf, rc); err != nil {
		log.WithError(err).Error("error reading bytes from bucket")
		return nil, err
	}

	log.WithField("filename", path).WithField("bucket", bucket).Info("upload to bucket complete")

	var printFileEntries []*pkg.PrintFileEntry
	err = json.Unmarshal(buf.Bytes(), &printFileEntries)
	if err != nil {
		log.WithError(err).Error("unable to marshall json payload - nacking message")
		return nil, err
	}
	printFile := pkg.PrintFile{
		PrintFiles: printFileEntries,
	}

	return &printFile, nil
}
