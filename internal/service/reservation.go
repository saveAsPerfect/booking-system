package service

import (
	"context"

	"github.com/saveAsPerfect/booking-system/internal/models"
	"github.com/saveAsPerfect/booking-system/internal/repository"
)

type ReservationService struct {
	repo repository.ReservationRepository
}

func NewReservationService(repo repository.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	// Здесь можно добавить дополнительную бизнес-логику, если необходимо
	return s.repo.CreateReservation(ctx, reservation)
}

func (s *ReservationService) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {
	// Здесь можно добавить дополнительную бизнес-логику, если необходимо
	return s.repo.GetReservations(ctx, roomID)
}