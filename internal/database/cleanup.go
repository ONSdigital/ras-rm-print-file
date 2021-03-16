package database

import (
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type CleanUp struct {
	store pkg.Store
}

func (c CleanUp) Start() {
	configDelay := viper.GetInt64("CLEANUP_DELAY")
	logger.Debug("retrieving clean up delay setting from config",
		zap.Int64("delay", configDelay))

	delay := time.Duration(configDelay) * time.Hour
	for {
		logger.Info("about to run clean up service")
		c.store = &DataStore{}
		c.process()
		logger.Info("sleeping clean up service",
			zap.String("delay", delay.String()))
		time.Sleep(delay)
	}
}

func (c CleanUp) process() {
	err := c.store.Init()
	defer c.store.Close()
	if err != nil {
		logger.Error("unable to initialise storage",
			zap.Error(err))
		return
	}
	printRequests, err := c.store.FindComplete()
	if err != nil {
		logger.Error("unable to find completed print file requests",
			zap.Error(err))
		return
	}
	if printRequests == nil {
		logger.Info("no completed print file requests to reprocess")
		return
	}
	complete := len(printRequests)
	logger.Info("found all completed print file requests",
		zap.Int("complete", complete))

	for _, printRequest := range printRequests {
		logger.Info("reprocessing print request")
		if printRequest.Status.Completed {
			now := time.Now()
			updated := printRequest.Updated
			duration := now.Sub(updated)
			retention := viper.GetInt64("CLEANUP_RETENTION")
			retentionDuration := time.Duration(retention) * time.Hour
			if duration >= retentionDuration {
				logger.Info("deleting print file request as its older than retention period",
					zap.Duration("retention", retentionDuration),
					zap.Duration("duration ", duration))
				c.store.Delete(printRequest)
			}
		} else {
			logger.Warn("unexpected incomplete print file request during clean up")
		}

	}
}
