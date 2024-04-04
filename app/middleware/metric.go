package middleware

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/DataDog/datadog-go/statsd"
)

type Metric struct{ Metric *statsd.Client }
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (m *Metric) Inc(name string) {
	err := m.Metric.Incr(
		name,
		[]string{},
		1,
	)
	if err != nil && err != statsd.ErrNoClient {
		log.Warnf("failed to incr metric: %s", err.Error())
	}
}

func (m *Metric) WithStatsd(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lrw := NewLoggingResponseWriter(w)
		h.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode
		m.Inc(fmt.Sprintf("%s.%s", r.Method, r.URL.Path[1:]))

		if statusCode != http.StatusOK {
			m.Inc(fmt.Sprintf("%s.%s.%d", r.Method, r.URL.Path[1:], statusCode))
		}
	})
}
