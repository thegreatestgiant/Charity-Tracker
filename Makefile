up:
	docker compose up -d
	

down:
	docker compose down

downF:
	docker compose down -v

logs:
	docker compose logs -f

run:
	go run backend/main.go
