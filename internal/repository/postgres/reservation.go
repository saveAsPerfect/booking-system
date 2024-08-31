package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saveAsPerfect/booking-system/internal/models"
	"github.com/saveAsPerfect/booking-system/internal/repository"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) repository.ReservationRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO reservations (room_id, start_time, end_time)
		VALUES ($1, $2, $3)
	`, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {

	rows, err := r.pool.Query(ctx, `
		SELECT id, room_id, start_time, end_time
		FROM reservations
		WHERE room_id = $1
		ORDER BY start_time
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(&r.ID, &r.RoomID, &r.StartTime, &r.EndTime)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *PostgresRepository) CheckReservation(ctx context.Context, reservation models.Reservation) error{

	var existingReservations []models.Reservation
	query := `
		SELECT *
		FROM reservations
		WHERE room_id = $1
		AND (
			(start_time BETWEEN $2 AND $3)
			OR (end_time BETWEEN $2 AND $3)
			OR (start_time < $2 AND end_time > $3)
		);
	`
	rows, err := r.pool.Query(ctx, query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var existingReservation models.Reservation
		err := rows.Scan(&existingReservation.ID, &existingReservation.RoomID, &existingReservation.StartTime, &existingReservation.EndTime)
		if err != nil {
			return err
		}
		existingReservations = append(existingReservations, existingReservation)
	}

	if len(existingReservations) > 0 {
		return errors.New("the room is reserved for this time")
	}
	return nil
}
