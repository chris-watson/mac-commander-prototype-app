package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chris-watson/mac-windows-installer-app/pkg/adapter/handler"
	"github.com/chris-watson/mac-windows-installer-app/pkg/adapter/handler/model"
	"github.com/chris-watson/mac-windows-installer-app/pkg/domain"
)

type mockCommanderService struct{}

var mockService = &mockCommanderService{}

func (m *mockCommanderService) HandlePing(host string) (domain.PingResult, error) {
	return domain.PingResult{
		Successful: true,
		Time:       100,
	}, nil
}

func (m *mockCommanderService) GetSystemInfo() (domain.SystemInfo, error) {
	return domain.SystemInfo{
		Hostname:  "test-host",
		IPAddress: "192.168.1.100",
	}, nil
}

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name           string
		request        model.CommandRequest
		expectedStatus int
		wantSuccess    bool
	}{
		{
			name: "ping google",
			request: model.CommandRequest{
				Type:    "ping",
				Payload: "google.com",
			},
			expectedStatus: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name: "ping localhost",
			request: model.CommandRequest{
				Type:    "ping",
				Payload: "localhost",
			},
			expectedStatus: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name: "get system info",
			request: model.CommandRequest{
				Type: "sysinfo",
			},
			expectedStatus: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name: "invalid command",
			request: model.CommandRequest{
				Type: "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			wantSuccess:    false,
		},
	}

	handler := handler.NewHandler(mockService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/command", bytes.NewBuffer(payload))
			w := httptest.NewRecorder()

			handler.HandleCommand(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response model.CommandResponse
			json.NewDecoder(w.Body).Decode(&response)
			if response.Success != tt.wantSuccess {
				t.Errorf("expected success %v, got %v", tt.wantSuccess, response.Success)
			}
		})
	}
}
