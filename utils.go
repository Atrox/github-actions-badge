package main

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func sendJSONResponse(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	resp := &response{}
	if err == nil {
		statusCode = http.StatusOK

		resp.Success = true
	} else {
		statusCode = http.StatusInternalServerError

		resp.Success = false
		resp.Message = err.Error()

		sentry.CaptureException(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}

func sendEndpointResponse(w http.ResponseWriter, r *http.Request, endpoint *Endpoint) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(endpoint)
	if err != nil {
		sendJSONResponse(w, r, err)
		return
	}
}
