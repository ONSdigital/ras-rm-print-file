package retry

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/database"
	"github.com/ONSdigital/ras-rm-print-file/internal/processor"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type BackoffRetry struct {
	store pkg.Store
	printer pkg.Printer
}

func(b BackoffRetry) Start()  {
	configDelay := viper.GetInt64("RETRY_DELAY")
	log.WithField("delay", configDelay).Debug("retrieving delay setting from config")

	delay := time.Duration(configDelay) * time.Millisecond
	b.printer = processor.Create()
	b.store = &database.DataStore{}
	for {
		log.Info("about to run retry service")
		b.process()
		log.WithField("delay", delay.String()).Info("retry sleeping")
		time.Sleep(delay)
		log.Info("retry sleep complete")
	}
}

func(b BackoffRetry) process() {
	err := b.store.Init()
	if err != nil {
		log.WithError(err).Error("unable to initialise storage")
		return
	}
	printRequests, err := b.store.FindIncomplete()
	if err != nil {
		log.WithError(err).Error("unable to find incomplete print file requests")
		return
	}
	if printRequests == nil {
		log.Info("no incomplete print file requests to reprocess")
		return
	}
	incomplete := len(printRequests)
	log.WithField("incomplete", incomplete).Info("finding all incomplete print file requests")

	for _, printRequest := range printRequests {
		log.Info("reprocessing print request")

		go b.printer.ReProcess(printRequest)
	}
}
