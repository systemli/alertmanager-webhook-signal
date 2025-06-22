package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (s *ConfigSuite) TestBuildConfig() {
	// Set environment variables for testing
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LISTEN_ADDR", ":9090")
	os.Setenv("SIGNAL_API_URL", "https://signal-cli.example.org/api/v1/rpc")
	os.Setenv("SIGNAL_API_USER", "user")
	os.Setenv("SIGNAL_API_PASSWORD", "password")
	os.Setenv("SIGNAL_ACCOUNT", "+1234567890")
	os.Setenv("SIGNAL_RCPT_GROUP_ID", "group-id")

	cfg := BuildConfig()

	s.Equal("debug", cfg.LogLevel)
	s.Equal(":9090", cfg.ListenAddr)
	s.Equal("https://signal-cli.example.org/api/v1/rpc", cfg.SignalApiUrl)
	s.Equal("user", cfg.SignalApiUser)
	s.Equal("password", cfg.SignalApiPassword)
	s.Equal("+1234567890", cfg.SignalAccount)
	s.Equal("group-id", cfg.SignalRcptGroupId)
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
