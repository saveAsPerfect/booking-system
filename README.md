# Conference Room Booking System

This is a simple REST API for a conference room booking system.

## Prerequisites

- Docker
- Docker Compose
- Go 1.21 

## API Endpoints

- `POST /reservations`: Create a new reservation
- `GET /reservations/{room_id}`: Get all reservations for a specific room

## How to run

```
make run
```

```
make migrate-up
```

## Request body
```json
{
    "room_id":"hel",
    "start_time":"2024-09-09T10:00:00Z",
    "end_time":"2024-09-09T11:00:00Z"
}
```