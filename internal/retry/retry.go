package retry

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/database"
	"strconv"
	"time"

	"github.com/ONSdigital/ras-rm-print-file/internal/processor"
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BackoffRetry struct {
	store   pkg.Store
	printer pkg.Printer
}

func (b BackoffRetry) Start() {
	retryDelay := viper.GetString("RETRY_DELAY")
	logger.Debug("retrieving delay setting from config",
		zap.String("retryDelay", retryDelay))
	configDelay, err := strconv.Atoi(retryDelay)
	if err != nil {
		logger.Error("unable to convert retry delay to int, defaulting to an hour",
			zap.String("retryDelay", retryDelay))
		// defaulting to an hour
		configDelay = 3600000
	}

	delay := time.Duration(configDelay) * time.Millisecond
	for {
		logger.Info("about to run retry service")
		b.printer = processor.Create()
		logger.Info("creating datastore connection")
		b.store = &database.DataStore{}
		b.process()
		logger.Info("retry sleeping",
			zap.String("delay", delay.String()))
		time.Sleep(delay)
		logger.Info("retry sleep complete")
	}
}

func (b BackoffRetry) process() {
	err := b.store.Init()
	defer b.store.Close()
	if err != nil {
		logger.Error("unable to initialise storage",
			zap.Error(err))
		return
	}
	printRequests, err := b.store.FindIncomplete()
	if err != nil {
		logger.Error("unable to find incomplete print file requests",
			zap.Error(err))
		return
	}
	if printRequests == nil {
		logger.Info("no incomplete print file requests to reprocess")
		return
	}
	incomplete := len(printRequests)
	logger.Info("finding all incomplete print file requests",
		zap.Int("incomplete", incomplete))

	for _, printRequest := range printRequests {
		logger.Info("reprocessing print request")
		b.printer.ReProcess(printRequest)
	}
}
