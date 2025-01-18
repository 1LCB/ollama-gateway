package controllers

import (
	"fmt"
	"net/http"
	m "ollamaGateway/metrics"
	"ollamaGateway/utils"
	"strings"
	"time"
)

var (
	rateLimitMap = utils.NewTemporaryMap[string, []time.Time](1 * time.Hour)
)

func IPAndAPIKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestIp := strings.Split(r.RemoteAddr, ":")[0]

		if !cfg.CanRequestFromIP(requestIp) {
			http.Error(w, "Unauthorized IP", http.StatusForbidden)
			logger.Warning("Request denied for " + r.RemoteAddr + ". Requests from this IP are not allowed")
			return
		}

		if apiKey := r.Header.Get(cfg.AuthHeaderName); apiKey == "" || !cfg.HasAPIKey(apiKey) {
			http.Error(w, "Unauthorized API Key", http.StatusUnauthorized)
			logger.Warning(r.RemoteAddr + " has provided an invalid API Key. Request denied")
			return
		}

		next(w, r)
	}
}

func MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Metrics.Enabled {
			next(w, r)
			return
		}

		crw := &utils.CustomResponseWriter{
			ResponseWriter: w,
		}

		now := time.Now()
		next(crw, r)
		elapsedTime := time.Since(now)

		m.RequestsDuration.WithLabelValues(r.Method, r.URL.Path).Observe(elapsedTime.Seconds())
		m.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", crw.StatusCode)).Inc()

		requestIp := strings.Split(r.RemoteAddr, ":")[0]
		m.RequestsByIP.WithLabelValues(requestIp, r.URL.Path, fmt.Sprintf("%d", crw.StatusCode)).Inc()

		statusCode := crw.StatusCode
		if statusCode >= 200 && statusCode < 300 {
			m.RequestsSuccess.WithLabelValues(r.Method, r.URL.Path).Inc()
		} else {
			m.RequestsFailure.WithLabelValues(r.Method, r.URL.Path).Inc()
		}

	}
}

func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.RateLimit.Enabled {
			next(w, r)
			return
		}

		requestIp := strings.Split(r.RemoteAddr, ":")[0]
		if isRateLimited(requestIp) {
			logger.Error(r.RemoteAddr + " Rate limit exceeded")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		recordRequest(requestIp)

		logger.Info("Forwarding request for " + r.RemoteAddr)
		next(w, r)
	}
}

func isRateLimited(ip string) bool {
	requests, _ := rateLimitMap.Get(ip)

	window := time.Now().Add(-time.Duration(cfg.RateLimit.TimeWindowSeconds) * time.Second)
	var validRequests []time.Time
	for _, t := range requests {
		if t.After(window) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= cfg.RateLimit.MaxRequests {
		return true
	}

	rateLimitMap.Set(ip, validRequests, 600*time.Second)
	return false
}

func recordRequest(ip string) {
	requests, _ := rateLimitMap.Get(ip)
	rateLimitMap.Set(ip, append(requests, time.Now()), 0)
}
