package web

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/ras-rm-print-file/internal/gcs"
	"github.com/ONSdigital/ras-rm-print-file/internal/processor"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("setting content type to application/json")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Print(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Debug("post request received processing")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("unable to read body")
			w.WriteHeader(http.StatusInternalServerError)
		}



		log.WithField("reqBody", string(reqBody)).Debug("body of request")
		vars := mux.Vars(r)
		filename := vars["filename"]
		if filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "missing filename")
		}
		log.WithField("filename", filename).Info("received request to print file")
		var printFileEntries []*pkg.PrintFileEntry

		err = json.Unmarshal(reqBody, &printFileEntries)
		if err != nil {
			log.WithError(err).Error("unable to marshall json payload")
			w.WriteHeader(http.StatusBadRequest)
		}
		printFile := pkg.PrintFile{
			PrintFiles: printFileEntries,
		}

		w.WriteHeader(http.StatusAccepted)
		resp, _ := json.Marshal(printFile)
		log.WithField("resp", string(resp)).Debug("about to process")
		uploadBodyToGCS(reqBody, filename)

		//spawn a process to process the print file

		go processor.CreateAndProcess(filename, &printFile)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Info("print - method not allowed")
		fmt.Fprintf(w, "Only POST methods are supported.")
	}
}

func uploadBodyToGCS(reqBody []byte, filename string) {
	copyOfBody := make([]byte, len(reqBody))
	copy(copyOfBody, reqBody)
	copyFilename := filename + ".json"
	gcsUpload := gcs.GCSUpload{}
	gcsUpload.Init()
	gcsUpload.UploadFile(copyFilename, copyOfBody)
}

func Alive(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Debug("alive OK")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"OK\"}")
	default:
		log.Debug("alive -method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only GET methods are supported")
	}
}

func Ready(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Debug("ready OK")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"READY\"}")
	default:
		log.Debug("ready - method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only GET methods are supported")
	}
}
