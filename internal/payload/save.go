package payload

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/gcs"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
)

type Payload struct {
	gcsUpload pkg.Upload
}

func Create() Payload {
	payload := Payload{}
	payload.gcsUpload = &gcs.GCSUpload{}
	return payload
}

// This function saves the payload (pre-templating) to the bucket for safe keeping and monitoring
func (p Payload) Save(filename string, payload []byte) {
	//rename file to .json
	payloadFileName := strings.TrimSuffix(filename, path.Ext(filename)) + ".json"
	log.WithField("payloadFileName", payloadFileName).Info("saving json payload to GCS")
	err := p.gcsUpload.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise gcs connection")
		return
	}
	defer p.gcsUpload.Close()
	err = p.gcsUpload.UploadFile(payloadFileName, payload)
	if err != nil {
		log.WithError(err).Error("unable to upload json payload")
		return
	}
	log.Info("saved json payload")
}
