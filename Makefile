run:
	go run ./cmd/api

migrate_cli:
	@read -p "Enter migration name: " name; \
	migrate create -seq -ext=.sql -dir=./migrations $$name

# it will migrate our schema if the dirty flag is set to true
migrate_force:
	migrate -path ./migrations -database "postgres://greenlight:1234@localhost/greenlight?sslmode=disable" force $(version_of_migration)

migrate_up:
	migrate -path ./migrations -database "postgres://greenlight:1234@localhost/greenlight?sslmode=disable" up

migrate_current_version:
	migrate -path ./migrations -database "postgres://greenlight:1234@localhost/greenlight" version

migrate_switch:
	migrate -path ./migrations -database "postgres://greenlight:1234@localhost/greenlight" goto $(version_of_migration)

