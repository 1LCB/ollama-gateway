package utils

import (
	"log"
	"ollamaGateway/config"
)

var (
	cfg    = config.GetConfig()
	logger *gatewayLogger
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

type gatewayLogger struct {
	enabled bool
}

func GetLogger() *gatewayLogger {
	if logger == nil {
		ReloadLogger()
	}
	return logger
}

func ReloadLogger() {
	if logger == nil{
		logger = &gatewayLogger{
			enabled: cfg.Logging,
		}
		return
	}
	logger.enabled = cfg.Logging
}

func (l *gatewayLogger) println(color, level, value string) {
	if !l.enabled {
		return
	}
	log.Println("[" + color + level + Reset + "] " + value)
}

func (l *gatewayLogger) Info(value string) {
	l.println(Green, "INFO", value)
}

func (l *gatewayLogger) Warning(value string) {
	l.println(Yellow, "WARNING", value)
}

func (l *gatewayLogger) Error(value string) {
	l.println(Red, "ERROR", value)
}

func (l *gatewayLogger) Debug(value string) {
	l.println(Blue, "DEBUG", value)
}

func (l *gatewayLogger) Response(value string) {
	l.println(Magenta, "RESPONSE", value)
}
