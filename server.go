package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	router            *chi.Mux
	signalApiUrl      string
	signalApiUser     string
	signalApiPassword string
	signalAccount     string
	signalRcptGroupId string
}

func NewServer(signalApiUrl, signalApiUser, signalApiPassword, signalAccount, signalRcptGroupId string) *Server {
	return &Server{
		router:            chi.NewRouter(),
		signalApiUrl:      signalApiUrl,
		signalApiUser:     signalApiUser,
		signalApiPassword: signalApiPassword,
		signalAccount:     signalAccount,
		signalRcptGroupId: signalRcptGroupId,
	}
}

func (s *Server) Start(listenAddr string) error {
	s.RegisterRoutes()

	return http.ListenAndServe(listenAddr, s.router)
}

func (s *Server) RegisterRoutes() {
	// s.router.Use(s.AuthMiddleware)
	s.router.Post("/alertmanager", s.handleAlertmanagerAlert)
}

func (s *Server) handleAlertmanagerAlert(w http.ResponseWriter, r *http.Request) {
	logger.Info("Alertmanager alert received")

	var alert AlertmanagerAlert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	messages := s.AlertToMessages(alert)
	for _, message := range messages {
		logger.With(
			zap.String("message", message),
		).Info("Sending alert message to Signal Group")

		err = s.sendSignal(message)
		if err != nil {
			logger.Error("Failed to send Signal message", zap.Error(err))
			http.Error(w, "Failed to send Signal message", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
