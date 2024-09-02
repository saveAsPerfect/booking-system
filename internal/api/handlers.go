package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/saveAsPerfect/booking-system/internal/models"
	"github.com/saveAsPerfect/booking-system/internal/service"
)

type Handler struct {
	service *service.ReservationService
}

func NewHandler(service *service.ReservationService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.CreateReservation(r.Context(), reservation)
	if err != nil {
		if errors.Is(err, models.ErrorRoomAlreadyReserved) {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetReservations(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	reservations, err := h.service.GetReservations(r.Context(), roomID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if reservations == nil {
		http.Error(w, "No reservations found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(reservations)
}
