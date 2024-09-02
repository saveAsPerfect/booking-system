# Conference Room Booking System

This is a simple REST API for a conference room booking system.

## Prerequisites

- Docker
- Docker Compose
- Go 1.21 

## Running the Application

To start the application, run:

```
make run
```

This will build and start the API and PostgreSQL database using Docker Compose.

## Running Tests

To run the tests, execute:

```
make test
```

## API Endpoints

- `POST /reservations`: Create a new reservation
- `GET /reservations/{room_id}`: Get all reservations for a specific room

## Development

For local development, you can run the PostgreSQL database using Docker and the API locally:

1. Start the database:
   ```
   docker-compose up db
   ```

2. Run the API:
   ```
   go run cmd/api/main.go
   ```

Make sure to set the `DATABASE_URL` environment variable correctly when running the API locally.

[![Go][go-shield]][go-url]