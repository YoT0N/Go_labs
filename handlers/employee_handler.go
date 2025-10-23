package handlers

import (
	"bus-enterprise-api/models"
	"bus-enterprise-api/storage"
	"encoding/json"
	"net/http"
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
	respondJSON(w, http.StatusOK, employees)
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

	// Перевірка чи існує вже працівник з таким ID
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
	// Перевірка чи існує працівник
	if _, exists := h.storage.GetEmployee(id); !exists {
		respondError(w, http.StatusNotFound, "Employee not found")
		return
	}

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Перевірка що ID в тілі запиту співпадає з ID в URL
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
	// Перевірка чи існує працівник
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
