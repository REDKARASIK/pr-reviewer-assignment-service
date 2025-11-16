DC = docker compose

db-up:
	$(DC) up -d db

migrate:
	$(DC) run --rm migrate

app-up:
	$(DC) up -d app

up: db-up migrate app-up

down:
	$(DC) down

logs:
	$(DC) logs -f app db
