package middleware

import (
	"net/http"
)

const (
	// API ключ для авторизації
	validAPIKey = "my-secret-api-key-12345"
	// Назва заголовка для API ключа
	apiKeyHeader = "X-API-Key"
)

// AuthMiddleware перевіряє наявність і правильність API ключа
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отримуємо API ключ з заголовка
		apiKey := r.Header.Get(apiKeyHeader)

		// Перевіряємо чи ключ присутній
		if apiKey == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "API key is missing. Please provide X-API-Key header"}`))
			return
		}

		// Перевіряємо чи ключ правильний
		if apiKey != validAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid API key"}`))
			return
		}

		// Якщо все ОК, передаємо виконання наступному обробнику
		next.ServeHTTP(w, r)
	})
}
