package main

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/database"
	"net/http"

	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	"github.com/ONSdigital/ras-rm-print-file/internal/gcpubsub"
	"github.com/ONSdigital/ras-rm-print-file/internal/processor"
	"github.com/ONSdigital/ras-rm-print-file/internal/retry"
	"github.com/ONSdigital/ras-rm-print-file/internal/web"
	logger "github.com/ONSdigital/ras-rm-print-file/logging"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func configure() {
	// config
	viper.AutomaticEnv()
	config.SetDefaults()
}

func startRetryService() {
	logger.Info("starting retry service")
	br := retry.BackoffRetry{}
	br.Start()
	logger.Info("started retry service")
}

func startCleanUpService() {
	logger.Info("starting clean up service")
	cu := database.CleanUp{}
	cu.Start()
	logger.Info("started clean up service")
}

func startPubSubListener() {
	logger.Info("starting gcpubsub listener")
	s := gcpubsub.Subscriber{
		Printer: processor.Create(),
	}
	s.Start()
	logger.Info("started gcpubsub listener")
}

func main() {
	configure()
	logger.Info("starting print-file")

	// configure the gorilla router
	r := mux.NewRouter()
	r.Use(web.Middleware)
	r.HandleFunc("/alive", web.Alive)
	r.HandleFunc("/ready", web.Ready)
	http.Handle("/", r)

	go startRetryService()
	go startPubSubListener()
	go startCleanUpService()

	logger.Info("started")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("service failed")
	}
}
