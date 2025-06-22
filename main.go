package main

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	logLevel := "info"
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel = envLogLevel
	}

	atomic := zap.NewAtomicLevel()
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	atomic.SetLevel(level)
	logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atomic,
	))
}

func main() {
	cfg := BuildConfig()
	logger.Info("Starting server", zap.String("listenAddr", cfg.ListenAddr))
	s := NewServer(cfg.SignalApiUrl, cfg.SignalApiUser, cfg.SignalApiPassword, cfg.SignalAccount, cfg.SignalRcptGroupId)
	if err := s.Start(cfg.ListenAddr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
