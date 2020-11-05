package web

import (
	"fmt"
	"net/http"

	logger "github.com/ONSdigital/ras-rm-print-file/logging"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("setting content type to application/json")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Alive(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		logger.Debug("alive OK")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"OK\"}")
	default:
		logger.Debug("alive -method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only GET methods are supported")
	}
}

func Ready(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		logger.Debug("ready OK")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{\"status\": \"READY\"}")
	default:
		logger.Debug("ready - method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only GET methods are supported")
	}
}
