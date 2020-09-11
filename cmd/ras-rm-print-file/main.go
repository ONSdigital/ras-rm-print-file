package main

import (
	"github.com/ONSdigital/ras-rm-print-file/internal/config"
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

	log.Info("started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
