package main

import (
	"bytes"
	"text/template"

	"go.uber.org/zap"
)

const alertTemplate = `{{ if eq .Alert.Status "firing" }}❗{{ else }}✅{{end }} Alert {{ .Alertname }} is {{ .Alert.Status }}

{{- if gt (len (.Alert.Labels)) 0 }}

Labels:
{{- range $key, $value := .Alert.Labels }}
- {{ $key }}: {{ $value }}
{{- end }}
{{- end }}

{{- if gt (len (.Alert.Annotations)) 0 }}

Annotations:
{{- range $key, $value := .Alert.Annotations }}
- {{ $key }}: {{ $value }}
{{- end }}
{{- end }}

{{ .Alert.GeneratorURL }}
`

// AlertmanagerAlert represents the structure of the Alertmanager webhook alert.
type AlertmanagerAlert struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []AMAlert         `json:"alerts"`
}

// AMAlert represents an alert within the Alertmanager alert.
type AMAlert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type AlertTemplateData struct {
	Alertname string
	Alert     AMAlert
}

func (s *Server) AlertToMessages(alert AlertmanagerAlert) []string {
	messages := make([]string, 0, len(alert.Alerts))
	for _, alert := range alert.Alerts {
		message := alertToMessage(alert)
		if message == "" {
			logger.Error("Failed to convert alert to message", zap.String("alert", alert.Fingerprint))
			continue
		}

		messages = append(messages, message)
	}

	return messages
}

func alertToMessage(alert AMAlert) string {
	alertName := alert.Labels["alertname"]
	alertmanagerTemplate, err := template.New("alertmanager").Parse(alertTemplate)
	if err != nil {
		logger.Error("Failed to parse alert template", zap.Error(err))
		return ""
	}

	var message bytes.Buffer
	err = alertmanagerTemplate.Execute(&message, AlertTemplateData{
		Alertname: alertName,
		Alert:     alert,
	})

	if err != nil {
		logger.Error("Failed to prepare alert message", zap.Error(err))
		return ""
	}

	return message.String()
}
