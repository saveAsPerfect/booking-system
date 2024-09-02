package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		pool, err = pgxpool.New(context.Background(), dsn)
		if err == nil {
			break
		}
		log.Printf("Unable to connect to database, retrying in 5 seconds: %v\n", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Unable to connect to database after 5 attempts: %v\n", err)
	}
	defer pool.Close()

	repo := postgres.NewPostgresRepository(pool)
	service := service.NewReservationService(repo)
	handler := api.NewHandler(service)
	router := api.SetupRouter(handler)

	quit := make(chan os.Signal, 1)

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("failed to start server: ", err)
		}
	}()

	log.Printf("Starting server on http://%s:%s", cfg.Server.Host, cfg.Server.Port)

	<-quit
	log.Println("stopping server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println("error on server shutting down:", err.Error())
	}

}
