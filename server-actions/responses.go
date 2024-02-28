package server_actions

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errBody := ErrorResponse{Error: msg}
	dat, err := json.Marshal(errBody)
	if err != nil {
		log.Printf("ErrorResponse encoding JSON: %v", err)
		return
	}
	w.Write(dat)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ErrorResponse encoding JSON: %v", err)
		return
	}
	w.Write(dat)
}
