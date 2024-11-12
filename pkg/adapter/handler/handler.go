package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chris-watson/mac-windows-installer-app/pkg/adapter/handler/model"
	"github.com/chris-watson/mac-windows-installer-app/pkg/domain"
)

type CommanderService interface {
	HandlePing(host string) (domain.PingResult, error)
	GetSystemInfo() (domain.SystemInfo, error)
}

type Handler struct {
	svc CommanderService
}

func NewHandler(svc CommanderService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) HandleCommand(w http.ResponseWriter, r *http.Request) {
	var req model.CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.CommandResponse{
			Success: false, Error: "Invalid request format",
		})
		return
	}

	if err := req.Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, model.CommandResponse{
			Success: false, Error: err.Error(),
		})
		return
	}

	switch req.Type {
	case "ping":
		var host string
		if req.Payload != "" {
			host = req.Payload
		}

		res, err := h.svc.HandlePing(host)
		if err != nil {
			writeJSON(w, http.StatusOK, model.CommandResponse{
				Success: false, Error: err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, model.CommandResponse{Success: true, Data: res})

	case "sysinfo":
		res, err := h.svc.GetSystemInfo()
		if err != nil {
			writeJSON(w, http.StatusOK, model.CommandResponse{
				Success: false, Error: err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, model.CommandResponse{Success: true, Data: res})

	default:
		writeJSON(w, http.StatusBadRequest, model.CommandResponse{
			Success: false, Error: "Unsupported command type",
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, response model.CommandResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("could not encode response: %v", err)
	}
}
