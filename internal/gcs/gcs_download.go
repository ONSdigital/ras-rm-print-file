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
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GCSDownload struct {
	ctx    context.Context
	client *storage.Client
}

func (d *GCSDownload) Init() error {
	var err error
	d.ctx = context.Background()
	logger.Info("initialising GCS connection")
	d.client, err = storage.NewClient(d.ctx)
	if err != nil {
		logger.Error("error creating gcp client", zap.Error(err))
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	logger.Info("connected to GCS")
	return nil
}

func (d *GCSDownload) Close() error {
	if d.client == nil {
		return errors.New("please initialise the connection")
	}
	logger.Info("closing connection to GCS")
	err := d.client.Close()
	logger.Info("GCS connection closed")
	return err
}

func (d *GCSDownload) DownloadFile(filename string) (*pkg.PrintFile, error) {
	// loads the payload from a GCS bucket
	if d.client == nil {
		return nil, errors.New("please initialise the connection")
	}
	bucket := viper.GetString("BUCKET_NAME")
	path := bucketPath(filename)

	logger.Info("downloading from bucket",
		zap.String("filename", path),
		zap.String("bucket", bucket))

	ctx, cancel := context.WithTimeout(d.ctx, time.Second*50)
	defer cancel()

	// GCSUpload an object with storage.Writer.
	rc, err := d.client.Bucket(bucket).Object(path).NewReader(ctx)
	if err != nil {
		logger.Error("error reading from bucket "+bucket+path, zap.Error(err))
		return nil, err
	}
	logger.Info("about to read contents from bucket",
		zap.String("filename", path),
		zap.String("bucket", bucket))

	buf := &bytes.Buffer{}
	defer rc.Close()
	if _, err := io.Copy(buf, rc); err != nil {
		logger.Error("error reading bytes from bucket", zap.Error(err))
		return nil, err
	}

	logger.Info("upload to bucket complete",
		zap.String("filename", path),
		zap.String("bucket", bucket))

	var printFileEntries []*pkg.PrintFileEntry
	err = json.Unmarshal(buf.Bytes(), &printFileEntries)
	if err != nil {
		logger.Error("unable to marshall json payload - nacking message", zap.Error(err))
		return nil, err
	}
	printFile := pkg.PrintFile{
		PrintFiles: printFileEntries,
	}

	return &printFile, nil
}
