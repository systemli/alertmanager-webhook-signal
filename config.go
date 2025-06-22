package main

import "os"

type Config struct {
	LogLevel          string
	ListenAddr        string
	SignalApiUrl      string
	SignalApiUser     string
	SignalApiPassword string
	SignalAccount     string
	SignalRcptGroupId string
}

func BuildConfig() *Config {
	cfg := &Config{
		LogLevel:   "info",
		ListenAddr: ":8080",
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if listenAddr := os.Getenv("LISTEN_ADDR"); listenAddr != "" {
		cfg.ListenAddr = listenAddr
	}
	cfg.SignalApiUrl = getEnvOrFatal("SIGNAL_API_URL")
	cfg.SignalAccount = getEnvOrFatal("SIGNAL_ACCOUNT")
	cfg.SignalRcptGroupId = getEnvOrFatal("SIGNAL_RCPT_GROUP_ID")
	cfg.SignalApiUser = getEnvOrFatal("SIGNAL_API_USER")
	cfg.SignalApiPassword = getEnvOrFatal("SIGNAL_API_PASSWORD")

	return cfg
}

func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Fatal("Environment variable " + key + " is required")
	}
	return val
}
