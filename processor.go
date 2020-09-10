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
	p.store.Init()
	pfr, err := p.store.Add(filename, printFile)
	if err != nil {
		log.WithError(err).Error("unable to DataStore print file request ")
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

	pfr.Status.UploadedGCS = p.uploadGCS(filename, buf)

	pfr.Status.UploadedSFTP = p.uploadSFTP(filename, buf)

	err = p.store.Update(pfr)
	if err != nil {
		log.WithError(err).Error("failed to Update database")
		//TODO set to not ready
		return err
	}
	return nil
}

func (p *Processor) uploadGCS(filename string, buffer *bytes.Buffer) bool {
	log.WithField("filename", filename).Info("uploading file to gcs")
	// first upload to GCS
	err := p.gcsUpload.Init()
	if err != nil {
		log.WithError(err).Error("failed to initialise GCS upload")
		return false
	}
	err = p.gcsUpload.UploadFile(filename, buffer.Bytes())
	if err != nil {
		log.WithError(err).Error("failed to upload to GCS")
		return false
	}
	return true
}

func (p *Processor) uploadSFTP(filename string, buffer *bytes.Buffer) bool {
	log.WithField("filename", filename).Info("uploading file to sftp")
	// and then to SFTP
	err := p.sftpUpload.Init()
	if err != nil {
		log.WithError(err).Error("failed to initialise SFTP upload")
		return false
	}
	err = p.sftpUpload.UploadFile(filename, buffer.Bytes())
	if err != nil {
		log.WithError(err).Error("failed to upload to SFTP")
		return false
	}
	return true
}
