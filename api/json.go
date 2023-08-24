package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	if statusCode > 499 {
		log.Println("Responed with 5xx error:", message)
	}

	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, statusCode, errResponse{
		Error: message,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	dataJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response. Payload: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dataJSON)
}
