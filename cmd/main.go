package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saveAsPerfect/booking-system/internal/api"
	"github.com/saveAsPerfect/booking-system/internal/config"
	"github.com/saveAsPerfect/booking-system/internal/repository/postgres"
	"github.com/saveAsPerfect/booking-system/internal/service"
)

func main() {

	cfg := config.MustLoad()
	
	dsn := fmt.Sprintf("postgres://%s:%v@%s:%v/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
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
	// TODD: add logger
	// TODO: gracefull shutdown
	// TODO: close db

}
