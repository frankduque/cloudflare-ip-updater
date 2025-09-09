package main

import (
	"github.com/frankduque/cloudflare-ip-updater/internal/config"
	"github.com/frankduque/cloudflare-ip-updater/internal/updater"
	"log"
	"os"
)

var (
	logFile *os.File
)

func main() {
	var err error
	logFile, err = os.OpenFile("cloudflare-updater.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	updater.Run(cfg)
}
