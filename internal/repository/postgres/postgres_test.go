// file: internal/repository/postgres/booking_test.go
package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/saveAsPerfect/booking-system/internal/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"github.com/saveAsPerfect/booking-system/internal/service"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
  //TODO: make config here
	dbPassword := "123"
	port := "5432"
	dbName := "bookings"


	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Printf("Could not connect to docker: %s", err)
		os.Exit(1)
	}

	resource, err := pool.Run("postgres", "13", []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
		fmt.Sprintf("POSTGRES_DB=%s", dbName),
	})
	if err != nil {
		fmt.Printf("Could not start resource: %s", err)
		os.Exit(1)
	}

	
	os.Setenv("TEST_DATABASE_URL", fmt.Sprintf("postgres://postgres:%s@localhost:%s/%s?sslmode=disable", 
		dbPassword, port,dbName))

	
	if err := pool.Retry(func() error {
		var err error
		db, err = pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	}); err != nil {
		fmt.Printf("Could not connect to docker: %s", err)
		os.Exit(1)
	}

	
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS reservations (
			id SERIAL PRIMARY KEY,
			room_id TEXT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		fmt.Printf("Could not create table: %s", err)
		os.Exit(1)
	}

	
	code := m.Run()


	db.Close()
	if err := pool.Purge(resource); err != nil {
		fmt.Printf("Could not purge resource: %s", err)
		os.Exit(1)
	}

	os.Exit(code)
}

func clearDatabase(t *testing.T) {
	t.Helper()
	_, err := db.Exec(context.Background(), "DELETE FROM reservations")
	if err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}
}

func TestSuccessfulBooking(t *testing.T) {
	clearDatabase(t)
	repo := NewPostgresRepository(db)
	ctx := context.Background()

	reservation := models.Reservation{
		RoomID:    "room1",
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	err := repo.CreateReservation(ctx, reservation)
	assert.NoError(t, err)

	reservations, err := repo.GetReservations(ctx, "room1")
	assert.NoError(t, err)
	assert.Len(t, reservations, 1)
}

func TestTimeConflictBooking(t *testing.T) {
	clearDatabase(t)
	repo := NewPostgresRepository(db)
	ctx := context.Background()

	reservation1 := models.Reservation{
		RoomID:    "room2",
		StartTime: time.Now().Add(time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	err := repo.CreateReservation(ctx, reservation1)
	assert.NoError(t, err)

	reservation2 := models.Reservation{
		RoomID:    "room2",
		StartTime: time.Now().Add(90 * time.Minute),
		EndTime:   time.Now().Add(150 * time.Minute),
	}

	err = repo.CheckReservation(ctx, reservation2)
	assert.Equal(t, models.ErrorRoomAlreadyReserved, err)
}

func TestConcurrentBooking(t *testing.T) {
	clearDatabase(t)
	
	repo := NewPostgresRepository(db)
	service := service.NewReservationService(repo)
	
	ctx := context.Background()

	reservation1 := models.Reservation{
		RoomID:    "room3",
		StartTime: time.Now().Add(3 * time.Hour),
		EndTime:   time.Now().Add(4 * time.Hour),
	}

	reservation2 := models.Reservation{
		RoomID:    "room3",
		StartTime: time.Now().Add(3 * time.Hour),
		EndTime:   time.Now().Add(4 * time.Hour),
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error

	go func() {
		defer wg.Done()
		err1 = service.CreateReservation(ctx, reservation1)
	}()

	go func() {
		defer wg.Done()
		err2 = service.CreateReservation(ctx, reservation2)
	}()

	wg.Wait()

	
	assert.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "Exactly one reservation should succeed")

	reservations, err := repo.GetReservations(ctx, "room3")
	assert.NoError(t, err)
	assert.Len(t, reservations, 1, "There should be exactly one reservation")
}

