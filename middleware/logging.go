package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// LoggingMiddleware логує всі HTTP запити у файл
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Логування запиту
		logMessage := fmt.Sprintf("[%s] %s %s %s\n",
			startTime.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
		)

		// Записуємо у файл
		file, err := os.OpenFile("logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Помилка відкриття файлу логів: %v", err)
		} else {
			defer file.Close()
			file.WriteString(logMessage)
		}

		// Виводимо також у консоль
		fmt.Print(logMessage)

		// Передаємо виконання наступному обробнику
		next.ServeHTTP(w, r)

		// Логуємо час виконання
		duration := time.Since(startTime)
		durationMessage := fmt.Sprintf("  └─ Completed in %v\n", duration)
		file.WriteString(durationMessage)
		fmt.Print(durationMessage)
	})
}
