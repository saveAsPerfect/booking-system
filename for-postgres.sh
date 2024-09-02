#!/bin/sh

docker run --name my-postgres -e POSTGRES_PASSWORD=123 -p 5432:5432 -d bookings