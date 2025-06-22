package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	server *Server
}

func (s *ServerSuite) SetupTest() {
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

func (s *ServerSuite) TestUnknownRoute() {
	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	s.server.router.ServeHTTP(w, req)

	s.Equal(404, w.Code)
}

func (s *ServerSuite) TestHandleAlertmanager() {
	s.Run("invalid request body", func() {
		req := httptest.NewRequest("POST", "/alertmanager", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()
		s.server.handleAlertmanagerAlert(w, req)
		s.Equal(400, w.Code)
	})

	s.Run("invalid signal API response", func() {
		alert := AlertmanagerAlert{
			GroupLabels:  map[string]string{"alertname": "TestAlert"},
			CommonLabels: map[string]string{"severity": "critical"},
			Alerts: []AMAlert{
				{
					Status:       "firing",
					Labels:       map[string]string{"alertname": "TestAlert", "severity": "critical"},
					Annotations:  map[string]string{"description": "Test alert firing"},
					GeneratorURL: "http://example.com/alert/12345",
				},
			},
		}
		jsonData, err := json.Marshal(alert)
		s.NoError(err)

		req := httptest.NewRequest("POST", "/alertmanager", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Mock Signal API response with invalid response
		gock.New("http://signal-cli.example.org").
			Reply(400)

		s.server.handleAlertmanagerAlert(w, req)
		s.True(gock.IsDone())
		s.Equal(400, w.Code)
	})

	s.Run("valid alert", func() {
		alert := AlertmanagerAlert{
			GroupLabels:  map[string]string{"alertname": "TestAlert"},
			CommonLabels: map[string]string{"severity": "critical"},
			Alerts: []AMAlert{
				{
					Status:       "firing",
					Labels:       map[string]string{"alertname": "TestAlert", "severity": "critical"},
					Annotations:  map[string]string{"description": "Test alert firing"},
					GeneratorURL: "http://example.com/alert/12345",
				},
			},
		}
		jsonData, err := json.Marshal(alert)
		s.NoError(err)

		req := httptest.NewRequest("POST", "/alertmanager", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Mock Signal API response with JSON-RPC response with a timestamp
		gock.New("http://signal-cli.example.org").
			Post("/api/v1/rpc").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      1,
				"result": map[string]interface{}{
					"Timestamp": 1,
				},
			})

		s.server.handleAlertmanagerAlert(w, req)
		s.True(gock.IsDone())
		s.Equal(200, w.Code)
	})
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
