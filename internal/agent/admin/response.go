package admin

import (
	"encoding/json"
	"net/http"
)

type respError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type response struct {
	Ok    bool       `json:"ok"`
	Data  any        `json:"data,omitempty"`
	Error *respError `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func RespondOK(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, response{Ok: true, Data: data})
}

func RespondError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, response{Ok: false, Error: &respError{Code: code, Message: message}})
}
