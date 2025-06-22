package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AlertSuite struct {
	suite.Suite
	server *Server
}

func (s *AlertSuite) SetupTest() {
	s.server = NewServer(
		"http://signal-cli.example.org/api/v1/rpc",
		"user",
		"password",
		"+1234567890",
		"group-id",
	)
}

func (s *AlertSuite) TestAlertToMessages() {
	s.Run("two valid alerts", func() {
		alert := AlertmanagerAlert{
			GroupLabels:  map[string]string{"alertname": "TestAlert"},
			CommonLabels: map[string]string{"severity": "critical"},
			Alerts: []AMAlert{
				{
					Status:       "firing",
					Labels:       map[string]string{},
					Annotations:  map[string]string{},
					GeneratorURL: "",
				},
				{
					Status:       "resolved",
					Labels:       map[string]string{},
					Annotations:  map[string]string{},
					GeneratorURL: "",
				},
			},
		}

		messages := s.server.AlertToMessages(alert)
		assert.Len(s.T(), messages, 2)
	})
}

func TestAlertToMessage(t *testing.T) {
	tests := []struct {
		name     string
		alert    AMAlert
		expected string
	}{
		{
			name: "Firing alert with labels and annotations",
			alert: AMAlert{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "HighCPUUsage",
					"severity":  "critical",
				},
				Annotations: map[string]string{
					"description": "CPU usage is above 90%",
					"summary":     "High CPU usage detected",
				},
				GeneratorURL: "http://example.com/alert/12345",
			},
			expected: "❗ Alert HighCPUUsage is firing\n\nLabels:\n- alertname: HighCPUUsage\n- severity: critical\n\nAnnotations:\n- description: CPU usage is above 90%\n- summary: High CPU usage detected\n\nhttp://example.com/alert/12345\n",
		},
		{
			name: "Resolved alert with no labels or annotations",
			alert: AMAlert{
				Status:       "resolved",
				Labels:       map[string]string{},
				Annotations:  map[string]string{},
				GeneratorURL: "http://example.com/alert/67890",
			},
			expected: "✅ Alert  is resolved\n\nhttp://example.com/alert/67890\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := alertToMessage(tt.alert)
			assert.Equal(t, tt.expected, message)
		})
	}
}
