package handlers

import (
	"encoding/json"
	"lab5/models"
	"lab5/storage"
	"net/http"
	"strconv"
	"strings"
)

type EmployeeHandler struct {
	storage *storage.Storage
}

func NewEmployeeHandler(storage *storage.Storage) *EmployeeHandler {
	return &EmployeeHandler{storage: storage}
}

func (h *EmployeeHandler) HandleEmployees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getEmployees(w, r)
	case http.MethodPost:
		h.createEmployee(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EmployeeHandler) HandleEmployee(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getEmployee(w, r, id)
	case http.MethodPut:
		h.updateEmployee(w, r, id)
	case http.MethodDelete:
		h.deleteEmployee(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EmployeeHandler) getEmployees(w http.ResponseWriter, r *http.Request) {
	employees := h.storage.GetAllEmployees()

	// Отримуємо параметри фільтрації з query string
	query := r.URL.Query()

	// Фільтруємо результати
	filtered := make([]models.Employee, 0)
	for _, emp := range employees {
		match := true

		// Фільтр по firstName (часткове співпадіння, регістронезалежне)
		if firstName := query.Get("firstName"); firstName != "" {
			if !strings.Contains(strings.ToLower(emp.PersonalInfo.FirstName), strings.ToLower(firstName)) {
				match = false
			}
		}

		// Фільтр по lastName (часткове співпадіння, регістронезалежне)
		if lastName := query.Get("lastName"); lastName != "" {
			if !strings.Contains(strings.ToLower(emp.PersonalInfo.LastName), strings.ToLower(lastName)) {
				match = false
			}
		}

		// Фільтр по position (часткове співпадіння, регістронезалежне)
		if position := query.Get("position"); position != "" {
			if !strings.Contains(strings.ToLower(emp.EmployeeData.Position), strings.ToLower(position)) {
				match = false
			}
		}

		// Фільтр по department (часткове співпадіння, регістронезалежне)
		if department := query.Get("department"); department != "" {
			if !strings.Contains(strings.ToLower(emp.EmployeeData.Department), strings.ToLower(department)) {
				match = false
			}
		}

		// Фільтр по мінімальному віку
		if minAge := query.Get("minAge"); minAge != "" {
			if age, err := strconv.Atoi(minAge); err == nil {
				if emp.PersonalInfo.Age < age {
					match = false
				}
			}
		}

		// Фільтр по максимальному віку
		if maxAge := query.Get("maxAge"); maxAge != "" {
			if age, err := strconv.Atoi(maxAge); err == nil {
				if emp.PersonalInfo.Age > age {
					match = false
				}
			}
		}

		// Фільтр по мінімальній зарплаті
		if minSalary := query.Get("minSalary"); minSalary != "" {
			if salary, err := strconv.ParseFloat(minSalary, 64); err == nil {
				if emp.EmployeeData.Salary < salary {
					match = false
				}
			}
		}

		// Фільтр по максимальній зарплаті
		if maxSalary := query.Get("maxSalary"); maxSalary != "" {
			if salary, err := strconv.ParseFloat(maxSalary, 64); err == nil {
				if emp.EmployeeData.Salary > salary {
					match = false
				}
			}
		}

		if match {
			filtered = append(filtered, emp)
		}
	}

	respondJSON(w, http.StatusOK, filtered)
}

func (h *EmployeeHandler) getEmployee(w http.ResponseWriter, r *http.Request, id string) {
	employee, exists := h.storage.GetEmployee(id)
	if !exists {
		respondError(w, http.StatusNotFound, "Employee not found")
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) createEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if employee.EmployeeID == "" {
		respondError(w, http.StatusBadRequest, "EmployeeID is required")
		return
	}

	if _, exists := h.storage.GetEmployee(employee.EmployeeID); exists {
		respondError(w, http.StatusConflict, "Employee with this ID already exists")
		return
	}

	if err := h.storage.CreateEmployee(employee); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create employee")
		return
	}

	respondJSON(w, http.StatusCreated, employee)
}

func (h *EmployeeHandler) updateEmployee(w http.ResponseWriter, r *http.Request, id string) {
	if _, exists := h.storage.GetEmployee(id); !exists {
		respondError(w, http.StatusNotFound, "Employee not found")
		return
	}

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if employee.EmployeeID != id {
		respondError(w, http.StatusBadRequest, "EmployeeID in body does not match URL ID")
		return
	}

	if err := h.storage.UpdateEmployee(id, employee); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update employee")
		return
	}

	respondJSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) deleteEmployee(w http.ResponseWriter, r *http.Request, id string) {
	if _, exists := h.storage.GetEmployee(id); !exists {
		respondError(w, http.StatusNotFound, "Employee not found")
		return
	}

	if err := h.storage.DeleteEmployee(id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete employee")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Employee deleted successfully"})
}
