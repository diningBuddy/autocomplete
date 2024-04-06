package handler

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	log.Errorf("responedError: %s", message)
	respondJSON(w, code, map[string]string{"error": message})
}
