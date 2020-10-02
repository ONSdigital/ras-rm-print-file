package web

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("setting content type to application/json")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
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
