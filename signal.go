package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

type authTransport struct {
	Username string
	Password string
	base     http.RoundTripper
}

type SignalGroupMessage struct {
	Account   string `json:"account"`
	GroupId   string `json:"group-id"`
	Message   string `json:"message"`
	TextStyle string `json:"text-style,omitempty"`
}

type SignalGroupResponse struct {
	Timestamp int `json:"timestamp"`
}

func (s *Server) sendSignal(message string) error {
	client := &http.Client{
		Transport: &authTransport{
			Username: s.signalApiUser,
			Password: s.signalApiPassword,
			base:     http.DefaultTransport,
		},
	}

	rpcClient := jsonrpc.NewClientWithOpts(s.signalApiUrl, &jsonrpc.RPCClientOpts{
		HTTPClient: client,
	})

	// Find positions for bold alert name
	pos := strings.Index(message, " Alert ") + len(" Alert ") - 2
	posEnd := strings.Index(message, " is ") - 2
	length := posEnd - pos
	logger.Info("Signal message", zap.Int("pos", pos), zap.Int("length", length))

	params := SignalGroupMessage{
		Account:   s.signalAccount,
		GroupId:   s.signalRcptGroupId,
		Message:   message,
		TextStyle: fmt.Sprintf("%d:%d:BOLD", pos, length),
	}

	var response SignalGroupResponse
	err := rpcClient.CallFor(context.Background(), &response, "send", &params)
	if err != nil {
		return err
	}
	if response.Timestamp == 0 {
		return errors.New("signal json-rpc response: no timestamp in response")
	}

	return nil
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(a.Username, a.Password)
	return a.base.RoundTrip(req)
}
