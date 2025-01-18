package main

import (
	"net/http"
	"ollamaGateway/config"

	c "ollamaGateway/controllers"
	"ollamaGateway/metrics"
	"ollamaGateway/utils"
)

func main() {
	if err := config.ReloadConfig(); err != nil {
		panic("It was not possible to reload the config file:" + err.Error())
	}
	utils.ReloadLogger()

	logger := utils.GetLogger()
	config := config.GetConfig()
	metrics.InitMetrics()

	middApplyer := utils.NewMiddlewareApplyer(
		c.RateLimitMiddleware,
		c.IPAndAPIKeyMiddleware,
		c.MetricsMiddleware,
	)

	http.HandleFunc("/", middApplyer.Apply(c.ProxyHandler))
	http.HandleFunc("GET /config/reload", middApplyer.Apply(c.ReloadConfigHandler))
	http.HandleFunc("GET /config/new-key", middApplyer.Apply(c.GenerateKeyHandler))
	http.HandleFunc("GET "+config.Metrics.Endpoint, c.MetricsHandler)

	logger.Info("Gateway running on " + config.GatewayAddress)
	if err := http.ListenAndServe(config.GatewayAddress, nil); err != nil {
		panic(err)
	}
}
