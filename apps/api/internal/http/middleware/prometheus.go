package middleware

import (
	"net/http"

	appmetrics "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Prometheus(metrics *appmetrics.Collector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(recorder, r)
			metrics.ObserveHTTPRequest(r.Method, r.URL.Path, recorder.statusCode)
		})
	}
}
