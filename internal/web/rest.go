package web

import (
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/ras-rm-print-file/internal/processor"
	"github.com/ONSdigital/ras-rm-print-file/pkg"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("setting content type to application/json")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Print(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		logrus.Debug("post request received processing")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.WithError(err).Error("unable to read body")
			w.WriteHeader(http.StatusInternalServerError)
		}
		logrus.WithField("reqBody", reqBody).Debug("body of request")
		vars := mux.Vars(r)
		filename := vars["filename"]
		if filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "missing filename")
		}
		logrus.WithField("filename", filename).Info("received request to print file")
		var printFileEntries []*pkg.PrintFileEntry
		err = json.Unmarshal(reqBody, &printFileEntries)
		if err != nil {
			logrus.WithError(err).Error("unable to marshall json payload")
			w.WriteHeader(http.StatusBadRequest)
		}
		printFile := pkg.PrintFile{
			PrintFiles: printFileEntries,
		}

		w.WriteHeader(http.StatusAccepted)
		resp, _ := json.Marshal(printFile)
		logrus.WithField("resp", string(resp)).Debug("about to process")
		//spawn a process to process the print file
		go processor.Process(filename, &printFile)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		logrus.Info("print - method not allowed")
		fmt.Fprintf(w, "Only POST methods are supported.")
	}
}

func Alive(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		logrus.Debug("alive OK")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"OK\"}")
	default:
		logrus.Debug("alive -method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only GET methods are supported")
	}
}

func Ready(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		logrus.Debug("ready OK")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"READY\"}")
	default:
		logrus.Debug("ready - method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only GET methods are supported")
	}
}



