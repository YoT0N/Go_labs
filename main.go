package main

import (
	"lab5/handlers"
	"lab5/storage"
	"fmt"
	"log"
	"net/http"
)

func main() {
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

	// Маршрутизація
	http.HandleFunc("/employees", employeeHandler.HandleEmployees)
	http.HandleFunc("/employees/", employeeHandler.HandleEmployee)
	http.HandleFunc("/buses", busHandler.HandleBuses)
	http.HandleFunc("/buses/", busHandler.HandleBus)
	http.HandleFunc("/routes", routeHandler.HandleRoutes)
	http.HandleFunc("/routes/", routeHandler.HandleRoute)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер запущено на порті %s\n", port)
	fmt.Println("Доступні ендпоінти:")
	fmt.Println("  GET/POST /employees")
	fmt.Println("  GET/PUT/DELETE /employees/{id}")
	fmt.Println("  GET/POST /buses")
	fmt.Println("  GET/PUT/DELETE /buses/{id}")
	fmt.Println("  GET/POST /routes")
	fmt.Println("  GET/PUT/DELETE /routes/{id}")

	log.Fatal(http.ListenAndServe(port, nil))
}
