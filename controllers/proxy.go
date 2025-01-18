package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"ollamaGateway/config"
	"ollamaGateway/utils"
	"sync"
)

var (
	cfg              = config.GetConfig()
	logger           = utils.GetLogger()
	currentIndex int = 0
	mu           sync.Mutex
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	address := getAddressRoundRobin()
	url, err := url.Parse(address)
	if err != nil {
		logger.Error("Error parsing Ollama address: " + err.Error())
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(url)
	reverseProxy.ServeHTTP(w, r)
}

func getAddressRoundRobin() string {
	mu.Lock()
	address := cfg.OllamaAddresses[currentIndex]
	if len(cfg.OllamaAddresses) > 1 {
		currentIndex = (currentIndex + 1) % len(cfg.OllamaAddresses)
	}
	mu.Unlock()
	return address
}
