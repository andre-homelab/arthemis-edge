up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f traefik

docs:
	swag init -g auth/main.go -o ./auth/docs