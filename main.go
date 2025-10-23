package main

import (
	"fmt"
	"lab5/handlers"
	"lab5/middleware"
	"lab5/storage"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Створюємо директорію для логів, якщо її немає
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Не вдалося створити директорію для логів: %v", err)
	}

	// Ініціалізація сховища
	storage := storage.NewStorage()

	// Завантаження даних
	if err := storage.LoadEmployees(); err != nil {
		log.Printf("Помилка завантаження працівників: %v", err)
	}
	if err := storage.LoadBuses(); err != nil {
		log.Printf("Помилка завантаження автобусів: %v", err)
	}
	if err := storage.LoadRoutes(); err != nil {
		log.Printf("Помилка завантаження маршрутів: %v", err)
	}

	// Ініціалізація обробників
	employeeHandler := handlers.NewEmployeeHandler(storage)
	busHandler := handlers.NewBusHandler(storage)
	routeHandler := handlers.NewRouteHandler(storage)

	// Створюємо мультиплексор
	mux := http.NewServeMux()

	// Маршрутизація
	mux.HandleFunc("/employees", employeeHandler.HandleEmployees)
	mux.HandleFunc("/employees/", employeeHandler.HandleEmployee)
	mux.HandleFunc("/buses", busHandler.HandleBuses)
	mux.HandleFunc("/buses/", busHandler.HandleBus)
	mux.HandleFunc("/routes", routeHandler.HandleRoutes)
	mux.HandleFunc("/routes/", routeHandler.HandleRoute)

	// Застосовуємо middleware (порядок важливий: спочатку логування, потім авторизація)
	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.AuthMiddleware(handler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер запущено на порті %s\n", port)
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("API ключ для авторизації: my-secret-api-key-12345")
	fmt.Println("Використовуйте заголовок: X-API-Key: my-secret-api-key-12345")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("\nДоступні ендпоінти:")
	fmt.Println("  GET/POST /employees")
	fmt.Println("  GET/PUT/DELETE /employees/{id}")
	fmt.Println("  GET/POST /buses")
	fmt.Println("  GET/PUT/DELETE /buses/{id}")
	fmt.Println("  GET/POST /routes")
	fmt.Println("  GET/PUT/DELETE /routes/{id}")
	fmt.Println("\nПриклади фільтрації (query parameters):")
	fmt.Println("  GET /employees?firstName=Іван&position=водій")
	fmt.Println("  GET /buses?minCapacity=30&maxYear=2021")
	fmt.Println("  GET /routes?startPoint=Центральний&maxDistance=20")
	fmt.Println("\nЛоги запитів зберігаються у файлі: logs/requests.log")
	fmt.Println("=" + strings.Repeat("=", 60) + "\n")

	log.Fatal(http.ListenAndServe(port, handler))
}
