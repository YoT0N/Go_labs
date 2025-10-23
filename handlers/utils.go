package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// відправляю JSON відповідь
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// відправляю помилку у форматі JSON
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// extractID витягує ID з URL шляху
// Наприклад: /employees/123 -> "123"
func extractID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}
