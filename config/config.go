package config

import (
	"encoding/json"
	"net"
	"ollamaGateway/metrics"
	"os"
)

var (
	config *Config
)

type RateLimit struct {
	MaxRequests       int  `json:"maxRequests"`
	TimeWindowSeconds int  `json:"timeWindowSeconds"`
	Enabled           bool `json:"enabled"`
}

type Security struct {
	AllowIps []string `json:"allowIps"`
	DenyIps  []string `json:"denyIps"`
}

type Metrics struct {
	Enabled   bool   `json:"enabled"`
	Endpoint  string `json:"endpoint"`
	Namespace string `json:"namespace"`
}

type Config struct {
	OllamaAddresses []string  `json:"ollamaAddresses"`
	GatewayAddress  string    `json:"gatewayAddress"`
	Logging         bool      `json:"logging"`
	AuthHeaderName  string    `json:"authHeaderName"`
	APIKeys         []string  `json:"apiKeys"`
	RateLimit       RateLimit `json:"rateLimit"`
	Security        Security  `json:"security"`
	Metrics         Metrics   `json:"metrics"`
}

func (c *Config) IsIPAllowed(ip string) bool {
	for _, allowedIP := range c.Security.AllowIps {
		if allowedIP == ip {
			return true
		}
	}
	return len(c.Security.AllowIps) == 0
}

func (c *Config) IsIPDenied(ip string) bool {
	for _, deniedIP := range c.Security.DenyIps {
		if deniedIP == ip {
			return true
		}
	}
	return false
}

func (c *Config) IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func (c *Config) HasAPIKey(apiKey string) bool {
	for _, ak := range c.APIKeys {
		if ak == apiKey {
			return true
		}
	}
	return false
}

func (c *Config) CanRequestFromIP(ip string) bool {
	if !c.IsValidIP(ip) {
		return false
	}

	if c.IsIPDenied(ip) {
		return false
	}

	return c.IsIPAllowed(ip)
}

func (c *Config) Update(updatedConfig *Config) {
	c.OllamaAddresses = updatedConfig.OllamaAddresses
	c.GatewayAddress = updatedConfig.GatewayAddress
	c.Logging = updatedConfig.Logging
	c.AuthHeaderName = updatedConfig.AuthHeaderName
	c.APIKeys = updatedConfig.APIKeys
	c.RateLimit = updatedConfig.RateLimit
	c.Security = updatedConfig.Security
	c.Metrics = updatedConfig.Metrics
}



func GetConfig() *Config {
	if config == nil {
		ReloadConfig()
	}
	return config
}

func ReloadConfig() error {
	configFileName := "./config.json"

	data, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	if config == nil {
		config = &cfg
	} else {
		config.Update(&cfg)
	}

	if config.Metrics.Enabled{
		metrics.APIKeysTotal.Set(float64(len(config.APIKeys)))
		metrics.OllamaServersTotal.Set(float64(len(config.OllamaAddresses)))
	}

	return nil
}
