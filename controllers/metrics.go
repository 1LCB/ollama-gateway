package controllers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)


func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if !cfg.Metrics.Enabled{
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Metrics are disabled"))
		return
	}
	promhttp.Handler().ServeHTTP(w, r)
}
