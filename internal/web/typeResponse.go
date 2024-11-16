package handlers

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, statusCode int, status, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
