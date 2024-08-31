package repository

import (
	"context"

	"github.com/saveAsPerfect/booking-system/internal/models"
)

type ReservationRepository interface {
	CreateReservation(ctx context.Context, reservation models.Reservation) error
	GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error)
	CheckReservation(ctx context.Context,reservation models.Reservation) (error)
}