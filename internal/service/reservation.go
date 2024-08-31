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
	
	if !isValidData(reservation) {
		return models.ErrorInvalidData
	}

	if err := s.repo.CheckReservation(ctx, reservation); err != nil {
		return err
	}
	return s.repo.CreateReservation(ctx, reservation)
}

func (s *ReservationService) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {
	return s.repo.GetReservations(ctx, roomID)
}

func (s *ReservationService) CheckReservation(ctx context.Context, reservation models.Reservation) error {
	return s.repo.CheckReservation(ctx, reservation)
}

func isValidData(reservation models.Reservation) bool {
	return reservation.RoomID != "" &&
		!reservation.StartTime.IsZero() &&
		!reservation.EndTime.IsZero()
}
