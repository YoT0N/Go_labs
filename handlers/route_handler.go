package handlers

import (
	"bus-enterprise-api/models"
	"bus-enterprise-api/storage"
	"encoding/json"
	"net/http"
)

type RouteHandler struct {
	storage *storage.Storage
}

func NewRouteHandler(storage *storage.Storage) *RouteHandler {
	return &RouteHandler{storage: storage}
}

func (h *RouteHandler) HandleRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getRoutes(w, r)
	case http.MethodPost:
		h.createRoute(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RouteHandler) HandleRoute(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid route ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getRoute(w, r, id)
	case http.MethodPut:
		h.updateRoute(w, r, id)
	case http.MethodDelete:
		h.deleteRoute(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RouteHandler) getRoutes(w http.ResponseWriter, r *http.Request) {
	routes := h.storage.GetAllRoutes()
	respondJSON(w, http.StatusOK, routes)
}

func (h *RouteHandler) getRoute(w http.ResponseWriter, r *http.Request, id string) {
	route, exists := h.storage.GetRoute(id)
	if !exists {
		respondError(w, http.StatusNotFound, "Route not found")
		return
	}
	respondJSON(w, http.StatusOK, route)
}

func (h *RouteHandler) createRoute(w http.ResponseWriter, r *http.Request) {
	var route models.Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if route.RouteID == "" {
		respondError(w, http.StatusBadRequest, "RouteID is required")
		return
	}

	if _, exists := h.storage.GetRoute(route.RouteID); exists {
		respondError(w, http.StatusConflict, "Route with this ID already exists")
		return
	}

	// Перевірка чи існують пов'язані сутності
	if route.BusAssigned != nil {
		if _, exists := h.storage.GetBus(route.BusAssigned.BusID); !exists {
			respondError(w, http.StatusBadRequest, "Assigned bus does not exist")
			return
		}
	}

	if route.Driver != nil {
		if _, exists := h.storage.GetEmployee(route.Driver.EmployeeID); !exists {
			respondError(w, http.StatusBadRequest, "Assigned driver does not exist")
			return
		}
	}

	if err := h.storage.CreateRoute(route); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create route")
		return
	}

	respondJSON(w, http.StatusCreated, route)
}

func (h *RouteHandler) updateRoute(w http.ResponseWriter, r *http.Request, id string) {
	if _, exists := h.storage.GetRoute(id); !exists {
		respondError(w, http.StatusNotFound, "Route not found")
		return
	}

	var route models.Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if route.RouteID != id {
		respondError(w, http.StatusBadRequest, "RouteID in body does not match URL ID")
		return
	}

	// Перевірка чи існують пов'язані сутності
	if route.BusAssigned != nil {
		if _, exists := h.storage.GetBus(route.BusAssigned.BusID); !exists {
			respondError(w, http.StatusBadRequest, "Assigned bus does not exist")
			return
		}
	}

	if route.Driver != nil {
		if _, exists := h.storage.GetEmployee(route.Driver.EmployeeID); !exists {
			respondError(w, http.StatusBadRequest, "Assigned driver does not exist")
			return
		}
	}

	if err := h.storage.UpdateRoute(id, route); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update route")
		return
	}

	respondJSON(w, http.StatusOK, route)
}

func (h *RouteHandler) deleteRoute(w http.ResponseWriter, r *http.Request, id string) {
	if _, exists := h.storage.GetRoute(id); !exists {
		respondError(w, http.StatusNotFound, "Route not found")
		return
	}

	if err := h.storage.DeleteRoute(id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete route")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Route deleted successfully"})
}
