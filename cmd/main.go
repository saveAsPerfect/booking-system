package main

import (
	"log"
	"net/http"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saveAsPerfect/booking-system/internal/api"
	"github.com/saveAsPerfect/booking-system/internal/config"
	"github.com/saveAsPerfect/booking-system/internal/repository/postgres"
	"github.com/saveAsPerfect/booking-system/internal/service"
)

func main() {

	cfg := config.MustLoad()

	connString := "postgres://postgres:123@localhost:5432/bookings?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	repo := postgres.NewPostgresRepository(pool)
	service := service.NewReservationService(repo)
	handler := api.NewHandler(service)
	router := api.SetupRouter(handler)

	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
