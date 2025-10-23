package handlers

import (
	"lab5/models"
	"lab5/storage"
	"encoding/json"
	"net/http"
)

type BusHandler struct {
	storage *storage.Storage
}

func NewBusHandler(storage *storage.Storage) *BusHandler {
	return &BusHandler{storage: storage}
}

func (h *BusHandler) HandleBuses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBuses(w, r)
	case http.MethodPost:
		h.createBus(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BusHandler) HandleBus(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid bus ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getBus(w, r, id)
	case http.MethodPut:
		h.updateBus(w, r, id)
	case http.MethodDelete:
		h.deleteBus(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BusHandler) getBuses(w http.ResponseWriter, r *http.Request) {
	buses := h.storage.GetAllBuses()
	respondJSON(w, http.StatusOK, buses)
}

func (h *BusHandler) getBus(w http.ResponseWriter, r *http.Request, id string) {
	bus, exists := h.storage.GetBus(id)
	if !exists {
		respondError(w, http.StatusNotFound, "Bus not found")
		return
	}
	respondJSON(w, http.StatusOK, bus)
}

func (h *BusHandler) createBus(w http.ResponseWriter, r *http.Request) {
	var bus models.Bus
	if err := json.NewDecoder(r.Body).Decode(&bus); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if bus.BusID == "" {
		respondError(w, http.StatusBadRequest, "BusID is required")
		return
	}

	if _, exists := h.storage.GetBus(bus.BusID); exists {
		respondError(w, http.StatusConflict, "Bus with this ID already exists")
		return
	}

	if err := h.storage.CreateBus(bus); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create bus")
		return
	}

	respondJSON(w, http.StatusCreated, bus)
}

func (h *BusHandler) updateBus(w http.ResponseWriter, r *http.Request, id string) {
	if _, exists := h.storage.GetBus(id); !exists {
		respondError(w, http.StatusNotFound, "Bus not found")
		return
	}

	var bus models.Bus
	if err := json.NewDecoder(r.Body).Decode(&bus); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if bus.BusID != id {
		respondError(w, http.StatusBadRequest, "BusID in body does not match URL ID")
		return
	}

	if err := h.storage.UpdateBus(id, bus); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update bus")
		return
	}

	respondJSON(w, http.StatusOK, bus)
}

func (h *BusHandler) deleteBus(w http.ResponseWriter, r *http.Request, id string) {
	if _, exists := h.storage.GetBus(id); !exists {
		respondError(w, http.StatusNotFound, "Bus not found")
		return
	}

	if err := h.storage.DeleteBus(id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete bus")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Bus deleted successfully"})
}
