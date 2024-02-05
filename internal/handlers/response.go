package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetErrorResponse(w http.ResponseWriter, handlerName string, err error, statusCode int) {
	w.WriteHeader(statusCode)
	buf := bytes.NewBufferString(handlerName)
	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteString("\n")
	_, _ = w.Write(buf.Bytes())
}

func GetSuccessResponseWithBody(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func GetErrorResponseWithBody(w http.ResponseWriter, statusCode int, body any) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		GetErrorResponse(w, "update", err, http.StatusInternalServerError)
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bodyData)
}