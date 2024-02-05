package handlers

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"net/http"
)

type UploadStateService interface {
	IsLastUploadSuccessful(ctx context.Context) (bool, error)
}

type StateHandler struct {
	uploadingInProcess *atomic.Bool
	uploadStateService UploadStateService
}

func NewStateHandler(uploadStateService UploadStateService, uploadingInProcess *atomic.Bool) *StateHandler {
	return &StateHandler{
		uploadStateService: uploadStateService,
		uploadingInProcess: uploadingInProcess,
	}
}

type StateResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

func (h *StateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stateResponse := &StateResponse{}

	if h.uploadingInProcess.Load() {
		stateResponse.Info = "updating"
	} else {
		isSuccessful, err := h.uploadStateService.IsLastUploadSuccessful(ctx)
		if err != nil {
			GetErrorResponse(w, "state", err, http.StatusInternalServerError)
		}
		if isSuccessful {
			stateResponse.Result = true
			stateResponse.Info = "ok"
		} else {
			stateResponse.Info = "empty"
		}
	}

	raw, err := json.Marshal(stateResponse)
	if err != nil {
		GetErrorResponse(w, "state", err, http.StatusInternalServerError)
	}

	GetSuccessResponseWithBody(w, raw)
}
