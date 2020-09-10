package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
)

type Processor struct {
	store Store
	gcsUpload Upload
	sftpUpload Upload
}


func Process(filename string, printFile *PrintFile) error {
	processor := &Processor{}
	processor.store = &DataStore{}
	processor.gcsUpload = &GCSUpload{}
	processor.sftpUpload = &SFTPUpload{}
	return processor.process(filename, printFile)
}

func (p *Processor) process(filename string, printFile *PrintFile) error {
	log.WithField("filename", filename).Info("processing print file")

	// first save the request to the DB
	err := p.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return err
	}
	pfr, err := p.store.Add(filename, printFile)
	if err != nil {
		log.WithError(err).Error("unable to store print file request ")
		return err
	}

	// first sanitise the data
	printFile.sanitise()

	// load the ApplyTemplate
	buf, err := printFile.ApplyTemplate(filename)
	if err != nil {
		return err
	}
	pfr.Status.TemplateComplete = true
	log.WithField("ApplyTemplate", printTemplate).WithField("filename", filename).Info("templating complete")

	// first upload to GCS
	pfr.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")

	// and then to SFTP
	pfr.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")

	err = p.store.Update(pfr)
	if err != nil {
		log.WithError(err).Error("failed to Update database")
		//TODO set to not ready
		return err
	}
	return nil
}

func upload(filename string, buffer *bytes.Buffer, uploader Upload, name string) bool {
	log.WithField("filename", filename).Infof("uploading file to %v", name)
	err := uploader.Init()
	if err != nil {
		log.WithError(err).Errorf("failed to initialise %v upload", name)
		return false
	}
	err = uploader.UploadFile(filename, buffer.Bytes())
	if err != nil {
		log.WithError(err).Errorf("failed to upload to %v", name)
		return false
	}
	return true
}

