package processor

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/ONSdigital/ras-rm-print-file/internal/database"
	"github.com/ONSdigital/ras-rm-print-file/internal/gcs"
	"github.com/ONSdigital/ras-rm-print-file/internal/sftp"
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"go.uber.org/zap"
)

var printTemplate = "printfile.tmpl"

type SDCPrinter struct {
	store       pkg.Store
	gcsUpload   pkg.Upload
	sftpUpload  pkg.Upload
	gcsDownload pkg.Download
}

func Create() *SDCPrinter {
	processor := &SDCPrinter{}
	processor.store = &database.DataStore{}
	processor.gcsUpload = &gcs.GCSUpload{}
	processor.sftpUpload = &sftp.SFTPUpload{}
	processor.gcsDownload = &gcs.GCSDownload{}
	return processor
}

func (p *SDCPrinter) Process(filename string, datafileName string) error {
	logger.Info("processing print file",
		zap.String("filename", filename))

	// first save the request to the DB
	err := p.store.Init()
	defer p.store.Close()
	if err != nil {
		logger.Error("unable to initialise storage",
			zap.Error(err))
		return err
	}
	printFileRequest, err := p.store.Add(filename, datafileName)
	if err != nil {
		logger.Error("unable to store print file request ",
			zap.Error(err))
		return err
	}
	// load the data file
	err = p.gcsDownload.Init()
	defer p.gcsDownload.Close()
	if err != nil {
		logger.Error("unable to initialise storage",
			zap.Error(err))
		return err
	}

	pf, err := p.gcsDownload.DownloadFile(datafileName)
	if err != nil {
		logger.Error("unable to load data file",
			zap.String("dataFile", datafileName))
		return err
	}

	// first sanitise the data
	sanitise(pf)

	// load the ApplyTemplate
	buf, err := applyTemplate(pf)
	if err != nil {
		return err
	}
	printFileRequest.Status.Templated = true
	logger.Info("templating complete",
		zap.String("ApplyTemplate", printTemplate),
		zap.String("filename", filename))

	// first upload to GCS
	printFileRequest.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")

	// and then to SFTP
	printFileRequest.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")

	// check if it's completed
	printFileRequest.Status.Completed = isComplete(printFileRequest)

	err = p.store.Update(printFileRequest)
	if err != nil {
		logger.Error("failed to Update database",
			zap.Error(err))
		return err
	}
	return nil
}

func (p *SDCPrinter) ReProcess(pfr *pkg.PrintFileRequest) error {
	filename := pfr.PrintFilename
	logger.Info("processing print file",
		zap.String("filename", filename))

	// increment the number of attempts
	numberOfAttempts := pfr.Attempts
	pfr.Attempts = numberOfAttempts + 1

	// load the data file
	err := p.gcsDownload.Init()
	if err != nil {
		logger.Error("unable to initialise storage",
			zap.Error(err))
		return err
	}
	defer p.gcsDownload.Close()

	printFile, err := p.gcsDownload.DownloadFile(pfr.DataFilename)
	if err != nil {
		logger.Error("unable to load data file",
			zap.String("dataFile", pfr.DataFilename))
		return err
	}

	// load the ApplyTemplate
	buf, err := applyTemplate(printFile)
	if err != nil {
		return err
	}
	pfr.Status.Templated = true
	logger.Info("templating complete",
		zap.String("ApplyTemplate", printTemplate),
		zap.String("filename", filename))

	// first upload to GCS
	if !pfr.Status.UploadedGCS {
		logger.Info("print request not uploaded to GCS retrying")
		pfr.Status.UploadedGCS = upload(filename, buf, p.gcsUpload, "gcs")
	}
	// and then to SFTP
	if !pfr.Status.UploadedSFTP {
		logger.Info("print request not uploaded to SFTP retrying")
		pfr.Status.UploadedSFTP = upload(filename, buf, p.sftpUpload, "sftp")
	}

	// check if it's completed
	pfr.Status.Completed = isComplete(pfr)

	// first save the request to the DB
	err = p.store.Init()
	if err != nil {
		logger.Error("unable to initialise storage",
			zap.Error(err))
		return err
	}
	defer p.store.Close()
	err = p.store.Update(pfr)
	if err != nil {
		logger.Error("failed to Update database",
			zap.Error(err))
		return err
	}
	return nil
}

func isComplete(printFileRequest *pkg.PrintFileRequest) bool {
	return printFileRequest.Status.Templated && printFileRequest.Status.UploadedGCS && printFileRequest.Status.UploadedSFTP
}

func upload(filename string, buffer *bytes.Buffer, uploader pkg.Upload, name string) bool {
	logger.Info("uploading file to ",
		zap.String("filename", filename))
	err := uploader.Init()
	if err != nil {
		logger.Error("failed to initialise upload to ",
			zap.String("name", name),
			zap.Error(err))
		return false
	}
	defer uploader.Close()
	err = uploader.UploadFile(filename, buffer.Bytes())
	if err != nil {
		logger.Error("failed to upload to ",
			zap.String("name", name),
			zap.Error(err))
		return false
	}
	return true
}

func sanitise(pf *pkg.PrintFile) {
	logger.Info("sanitising print file to match expected outcomes")
	for _, pfe := range pf.PrintFiles {
		pfe.SampleUnitRef = strings.TrimSpace(pfe.SampleUnitRef)
		pfe.Iac = nullIfEmpty(strings.TrimSpace(pfe.Iac))
		pfe.CaseGroupStatus = nullIfEmpty(pfe.CaseGroupStatus)
		pfe.EnrolmentStatus = nullIfEmpty(pfe.EnrolmentStatus)
		pfe.RespondentStatus = nullIfEmpty(pfe.RespondentStatus)
		pfe.Contact.Forename = nullIfEmpty(pfe.Contact.Forename)
		pfe.Contact.Surname = nullIfEmpty(pfe.Contact.Surname)
		pfe.Contact.EmailAddress = nullIfEmpty(pfe.Contact.EmailAddress)
		pfe.Region = nullIfEmpty(pfe.Region)
	}
}

func nullIfEmpty(value string) string {
	if value == "" {
		return "null"
	}
	return value
}

func applyTemplate(pf *pkg.PrintFile) (*bytes.Buffer, error) {
	// find the template
	wd, err := os.Getwd()
	if err != nil {
		logger.Error("unable to load template",
			zap.Error(err))
	}
	// Template location
	templateLocation := wd + "/templates/" + printTemplate

	// load the ApplyTemplate
	logger.Info("about to load ApplyTemplate",
		zap.String("ApplyTemplate", printTemplate))
	t, err := template.New(printTemplate).ParseFiles(templateLocation)
	if err != nil {
		logger.Error("failed to find ApplyTemplate",
			zap.Error(err))
		return nil, err
	}

	logger.Info("about to process ApplyTemplate",
		zap.String("ApplyTemplate", printTemplate))
	// create a bytes buffer and run the ApplyTemplate engine
	buf := &bytes.Buffer{}
	err = t.Execute(buf, pf)
	if err != nil {
		logger.Error("failed to process ApplyTemplate",
			zap.Error(err))
		return nil, err
	}
	logger.Info("templating complete",
		zap.String("ApplyTemplate", printTemplate))
	return buf, nil
}
