package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"ollamaGateway/config"
	"ollamaGateway/utils"
)

var (
	cfg          = config.GetConfig()
	logger       = utils.GetLogger()
	loadBalancer = utils.NewLoadBalancer(cfg.OllamaAddresses)
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	address, err := loadBalancer.GetServerByRoundRobin()
	if err != nil {
		logger.Error("Healthy Backend Servers are not available")
		return
	}

	logger.Info("Server selected: " + address)
	url, err := url.Parse(address)
	if err != nil {
		logger.Error("Error parsing Ollama address: " + err.Error())
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	reverseProxy.ServeHTTP(w, r)
}
