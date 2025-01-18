package controllers

import (
	"encoding/json"
	"net/http"
	"ollamaGateway/config"
	"ollamaGateway/utils"
	"os"
)

func GenerateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := utils.GenerateKey(32)

	cfg.APIKeys = append(cfg.APIKeys, key)

	file, err := os.OpenFile("./config.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(cfg)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}

func ReloadConfigHandler(w http.ResponseWriter, r *http.Request) {
	if err := config.ReloadConfig(); err != nil {
		http.Error(w, "Failed to reload the config file: "+err.Error(), http.StatusInternalServerError)
		logger.Error("Failed to reload the config file: " + err.Error())
		return
	}
	utils.ReloadLogger()
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	logger.Info("Config and Logger reloaded successfully!")
}
