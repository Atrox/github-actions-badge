package main

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func sendJSONResponse(w http.ResponseWriter, err error) {
	var statusCode int
	resp := response{}
	if err == nil {
		statusCode = http.StatusOK

		resp.Success = true
	} else {
		statusCode = http.StatusInternalServerError

		resp.Success = false
		resp.Message = err.Error()

		// raven.CaptureError(err, nil)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}
