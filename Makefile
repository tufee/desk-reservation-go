DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
MIGRATION_DIR=./internal/infra/migration

migration:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ You need to pass migration name: make migration NAME=migration_name"; \
		exit 1; \
		fi
	migrate create -ext sql -dir $(MIGRATION_DIR) $(NAME)

up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down

drop:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" drop -f
