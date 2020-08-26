package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func print(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Error("unable to read body")
				w.WriteHeader(http.StatusInternalServerError)
			}
			var printFileEntries []*PrintFileEntry
			err = json.Unmarshal(reqBody, &printFileEntries)
			if err != nil {
				log.WithError(err).Error("unable to marshall json payload")
				w.WriteHeader(http.StatusBadRequest)
			}
			printFile := PrintFile{
				PrintFiles: printFileEntries,
			}

			//spawn a process to process the printfile
			go printFile.process()
			w.WriteHeader(http.StatusAccepted)
			resp, _ := json.Marshal(printFile)

			fmt.Fprintln(w, string(resp))

		default:
			fmt.Fprintf(w, "Only POST methods are supported.")
	}
}

func configureLogging() {
	verbose := viper.GetBool("VERBOSE")
	if verbose {
		//anything debug and above
		log.SetLevel(log.DebugLevel)
	} else {
		//otherwise keep it to info
		log.SetLevel(log.InfoLevel)
	}
}

func setDefaults() {
	viper.SetDefault("VERBOSE", true)
	viper.SetDefault("BUCKET_NAME", "ras-rm-printfile")
}

func configure() {
	//config
	viper.AutomaticEnv()
	setDefaults()
	configureLogging()
}


func main() {
	configure()
	http.HandleFunc("/print", print)
	log.Fatal(http.ListenAndServe(":8080", nil))
}