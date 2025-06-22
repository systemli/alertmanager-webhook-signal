package main

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type SignalSuite struct {
	suite.Suite
	server *Server
}

func (s *SignalSuite) SetupTest() {
	gock.DisableNetworking()
	defer gock.Off()

	s.server = NewServer(
		"http://signal-cli.example.org/api/v1/rpc",
		"user",
		"password",
		"+1234567890",
		"group-id",
	)
}

func (s *SignalSuite) TestSendSignal() {
	s.Run("invalid signal API response", func() {
		// Mock Signal API response with invalid response
		gock.New("http://signal-cli.example.org").
			Reply(500)

		err := s.server.sendSignal("message")
		s.Error(err)
		s.True(gock.IsDone())
	})

	s.Run("no timestamp in response", func() {
		// Mock Signal API response with no timestamp
		gock.New("http://signal-cli.example.org").
			Post("/api/v1/rpc").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]interface{}{})

		err := s.server.sendSignal("message")
		s.Error(err)
		s.Contains(err.Error(), "signal json-rpc response: no timestamp in response")
	})

	s.Run("valid signal API response", func() {
		// Mock Signal API response with valid response
		gock.New("http://signal-cli.example.org").
			Post("/api/v1/rpc").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]interface{}{
				"timestamp": 1234567890,
			})
	})
}

func TestSignalSuite(t *testing.T) {
	suite.Run(t, new(SignalSuite))
}
