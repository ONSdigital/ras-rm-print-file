package main

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
	"github.com/ONSdigital/ras-rm-print-file/internal/retry"
	"github.com/ONSdigital/ras-rm-print-file/internal/web"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func configure() {
	//config
	viper.AutomaticEnv()
	config.SetDefaults()
	config.ConfigureLogging()
}

func startRetryService() {
	log.Info("starting retry service")
	br := retry.BackoffRetry{}
	go br.Start()
	log.Info("started retry service")
}

func main() {
	configure()
	log.Info("starting ras-rm-print-file")

	//configure the gorilla router
	r := mux.NewRouter()
	r.Use(web.Middleware)
	r.HandleFunc("/print/{filename}", web.Print)
	r.HandleFunc("/alive", web.Alive)
	r.HandleFunc("/ready", web.Ready)
	http.Handle("/", r)

	startRetryService()

	log.Info("started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
