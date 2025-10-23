package storage

import (
	"encoding/json"
	"lab5/models"
	"os"
	"sync"
)

type Storage struct {
	employees map[string]models.Employee
	buses     map[string]models.Bus
	routes    map[string]models.Route
	mutex     sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		employees: make(map[string]models.Employee),
		buses:     make(map[string]models.Bus),
		routes:    make(map[string]models.Route),
	}
}

func (s *Storage) loadFromFile(filename string, target interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(target)
}

func (s *Storage) saveToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Внутрішній метод без блокування (для використання всередині)
func (s *Storage) saveEmployeesInternal() error {
	var employees []models.Employee
	for _, employee := range s.employees {
		employees = append(employees, employee)
	}
	return s.saveToFile("storage/employees.json", employees)
}

// Методи для працівників
func (s *Storage) LoadEmployees() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var employees []models.Employee
	err := s.loadFromFile("storage/employees.json", &employees)
	if err != nil {
		return err
	}

	for _, employee := range employees {
		s.employees[employee.EmployeeID] = employee
	}
	return nil
}

func (s *Storage) SaveEmployees() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.saveEmployeesInternal()
}

func (s *Storage) CreateEmployee(employee models.Employee) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.employees[employee.EmployeeID] = employee
	return s.saveEmployeesInternal()
}

func (s *Storage) GetEmployee(id string) (models.Employee, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	employee, exists := s.employees[id]
	return employee, exists
}

func (s *Storage) GetAllEmployees() []models.Employee {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var employees []models.Employee
	for _, employee := range s.employees {
		employees = append(employees, employee)
	}
	return employees
}

func (s *Storage) UpdateEmployee(id string, employee models.Employee) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.employees[id] = employee
	return s.saveEmployeesInternal()
}

func (s *Storage) DeleteEmployee(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.employees, id)
	return s.saveEmployeesInternal()
}

// Внутрішній метод для автобусів
func (s *Storage) saveBusesInternal() error {
	var buses []models.Bus
	for _, bus := range s.buses {
		buses = append(buses, bus)
	}
	return s.saveToFile("storage/buses.json", buses)
}

// Методи для автобусів
func (s *Storage) LoadBuses() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var buses []models.Bus
	err := s.loadFromFile("storage/buses.json", &buses)
	if err != nil {
		return err
	}

	for _, bus := range buses {
		s.buses[bus.BusID] = bus
	}
	return nil
}

func (s *Storage) SaveBuses() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.saveBusesInternal()
}

func (s *Storage) CreateBus(bus models.Bus) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.buses[bus.BusID] = bus
	return s.saveBusesInternal()
}

func (s *Storage) GetBus(id string) (models.Bus, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	bus, exists := s.buses[id]
	return bus, exists
}

func (s *Storage) GetAllBuses() []models.Bus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var buses []models.Bus
	for _, bus := range s.buses {
		buses = append(buses, bus)
	}
	return buses
}

func (s *Storage) UpdateBus(id string, bus models.Bus) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.buses[id] = bus
	return s.saveBusesInternal()
}

func (s *Storage) DeleteBus(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.buses, id)
	return s.saveBusesInternal()
}

// Внутрішній метод для маршрутів
func (s *Storage) saveRoutesInternal() error {
	var routes []models.Route
	for _, route := range s.routes {
		routes = append(routes, route)
	}
	return s.saveToFile("storage/routes.json", routes)
}

// Методи для маршрутів
func (s *Storage) LoadRoutes() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var routes []models.Route
	err := s.loadFromFile("storage/routes.json", &routes)
	if err != nil {
		return err
	}

	for _, route := range routes {
		s.routes[route.RouteID] = route
	}
	return nil
}

func (s *Storage) SaveRoutes() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.saveRoutesInternal()
}

func (s *Storage) CreateRoute(route models.Route) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.routes[route.RouteID] = route
	return s.saveRoutesInternal()
}

func (s *Storage) GetRoute(id string) (models.Route, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	route, exists := s.routes[id]
	return route, exists
}

func (s *Storage) GetAllRoutes() []models.Route {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var routes []models.Route
	for _, route := range s.routes {
		routes = append(routes, route)
	}
	return routes
}

func (s *Storage) UpdateRoute(id string, route models.Route) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.routes[id] = route
	return s.saveRoutesInternal()
}

func (s *Storage) DeleteRoute(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.routes, id)
	return s.saveRoutesInternal()
}
