package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func print(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Error("unable to read body")
				w.WriteHeader(http.StatusInternalServerError)
			}
			vars := mux.Vars(r)
			filename := vars["filename"]
			if filename == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Missing filename")
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
			go printFile.process(filename)
			w.WriteHeader(http.StatusAccepted)
			resp, _ := json.Marshal(printFile)

			fmt.Fprintln(w, string(resp))

		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Only POST methods are supported.")
	}
}

func alive(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"OK\"}")
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Only GET methods are supporteds")
	}
}

func ready(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"READY\"}")
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Only GET methods are supporteds")
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
	viper.SetDefault("SFTP_HOST", "localhost")
	viper.SetDefault("SFTP_PORT", "22")
	viper.SetDefault("SFTP_USERNAME", "sftp")
	viper.SetDefault("SFTP_PASSWORD", "sftp")
	viper.SetDefault("SFTP_DIRECTORY", "printfiles")
}

func configure() {
	//config
	viper.AutomaticEnv()
	setDefaults()
	configureLogging()
}


func main() {
	configure()

	r := mux.NewRouter()
	r.Use(middleware)
	r.HandleFunc("/print/{filename}", print)
	r.HandleFunc("/alive", alive)
	r.HandleFunc("/ready", ready)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
