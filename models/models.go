package models

// Особиста інформація працівника
type PersonalInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Phone     string `json:"phone"`
}

// Облікові дані працівника
type EmployeeData struct {
	Position   string  `json:"position"`
	Salary     float64 `json:"salary"`
	HireDate   string  `json:"hireDate"`
	Department string  `json:"department"`
}

// Працівник
type Employee struct {
	PersonalInfo `json:"personalInfo"`
	EmployeeData `json:"employeeData"`
	EmployeeID   string `json:"employeeId"`
}

// Автобус
type Bus struct {
	BusID        string `json:"busId"`
	Model        string `json:"model"`
	Capacity     int    `json:"capacity"`
	Year         int    `json:"year"`
	LicensePlate string `json:"licensePlate"`
	Mileage      int    `json:"mileage"`
}

// Маршрут
type Route struct {
	RouteID     string    `json:"routeId"`
	RouteNumber string    `json:"routeNumber"`
	StartPoint  string    `json:"startPoint"`
	EndPoint    string    `json:"endPoint"`
	Distance    float64   `json:"distance"`
	Duration    float64   `json:"duration"` // у годинах
	BusAssigned *Bus      `json:"busAssigned"` // посилання на автобус
	Driver      *Employee `json:"driver"` // посилання на водія
}
