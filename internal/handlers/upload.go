package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type UploadService interface {
	Upload(ctx context.Context) error
}

type UploadHandler struct {
	uploadService      UploadService
	uploadingInProcess *atomic.Bool
}

func NewUploadHandler(uploadService UploadService, uploadingInProcess *atomic.Bool) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		uploadingInProcess: uploadingInProcess,
	}
}

type UploadResponse struct {
	Result bool   `json:"result"`
	Info   string `json:"info"`
	Code   int    `json:"code"`
}

func (h *UploadHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uploadResponse := &UploadResponse{}

	if h.uploadingInProcess.Load() {
		uploadResponse.Info = "updating"
		uploadResponse.Code = http.StatusInternalServerError
		GetErrorResponseWithBody(w, http.StatusInternalServerError, uploadResponse)
		return
	}

	h.uploadingInProcess.Store(true)
	err := h.uploadService.Upload(ctx)
	h.uploadingInProcess.Store(false)

	if err != nil {
		uploadResponse.Info = "service unavailable"
		uploadResponse.Code = http.StatusServiceUnavailable
		GetErrorResponseWithBody(w, http.StatusServiceUnavailable, uploadResponse)
		return
	}

	uploadResponse.Result = true
	uploadResponse.Code = 200

	raw, err := json.Marshal(uploadResponse)
	if err != nil {
		GetErrorResponse(w, "update", err, http.StatusInternalServerError)
	}

	GetSuccessResponseWithBody(w, raw)
}
