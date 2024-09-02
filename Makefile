
run:
	docker compose up --build -d

migrate-up:
	docker compose run migrate up 


clean:
	docker-compose down --volumes --remove-orphans
	docker system prune -f