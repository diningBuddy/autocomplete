package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Debugf("%s %s in %v", method, uri, duration)
	}
	return http.HandlerFunc(logFn)
}
