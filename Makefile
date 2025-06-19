#!make
include ./app.env
export $(shell sed 's/=.*//' app.env)
ROLLBACK_COUNT ?= 1

release: ## Release version for api
	@echo "Building release version..."
	@{ \
		current_full_version=$$(cat .build-version); \
		current_version=$$(echo $$current_full_version | cut -d'-' -f1); \
		IFS='.' read -r major minor patch <<< "$$current_version"; \
		patch=$$((10#$${patch})); \
		minor=$$((10#$${minor})); \
		major=$$((10#$${major})); \
		if [ "$$patch" -eq 999 ]; then \
			patch=0; \
			minor=$$((minor + 1)); \
			if [ "$$minor" -eq 10 ]; then \
				minor=0; \
				major=$$((major + 1)); \
			fi; \
		else \
			patch=$$((patch + 1)); \
		fi; \
		new_version="$${major}.$${minor}.$$(printf "%03d" $$patch)"; \
		git_hash=$$(git rev-parse --short=8 HEAD); \
		final_version="$${new_version}-g$${git_hash}"; \
		echo "$$final_version" > .build-version; \
		git add .build-version; \
		git commit -m "chore(release): build version $$final_version"; \
		echo "✅ Released version: $$final_version"; \
	}

build-image: ## Build image with current version
	@version=$$(cat .build-version); \
	echo "Building Docker image with version $$version..."; \
	docker build -t example-api:$$version .

run-image: ## Run image with current version
	@version=$$(cat .build-version); \
	echo "Running Docker image with version $$version..."; \
	docker run -d -p 8080:8080 example-api:$$version

push-image: ## Push image to registry
	@version=$$(cat .build-version); \
	echo "Pushing Docker image with version $$version..."; \
	docker push your-registry/example-api:$$version

# sqlc commands migration
create-migration: ## Create a new migration "Ex: make create-migration MIGRATION_NAME=create_customers"
	bash ./scripts/create-migration.sh $$MIGRATION_NAME

# sqlc commands database
create-version: # support create a new migration file with the version
	bash ./scripts/temp-db-migration.sh $$VERSION

migrate: # support migrate the migration file
	goose -dir db/migrations $$DATABASE_DRIVER "$$DATABASE_DRIVER://$$DATABASE_USERNAME:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_NAME?sslmode=disable" up

db-rollback-version: # support rollback the last migration version
	@echo "Rolling back and extracting the last migration version...";
	ROLLBACK_VERSION=$$(make db-rollback-one 2>&1 | grep -oE "migration [0-9]+" | awk '{print $$2}'); \
	if [ -z "$$ROLLBACK_VERSION" ]; then \
		echo "No migration version found, only rolling back one migration."; \
	else \
		echo "Found ROLLBACK_VERSION: $$ROLLBACK_VERSION"; \
		make create-version VERSION=$$ROLLBACK_VERSION; \
		make db-rollback-one; \
	fi

remove-migration-temp: # support remove the migration file temp from rollback
	echo "Finding migrations...temp from rollback..."
	files=$$(grep -rl "Migration missing from auto-generated" db/migrations); \
	if [ -n "$$files" ]; then \
		echo "Found files: $$files"; \
		echo "$$files" | xargs rm; \
	else \
		echo "No matching migrations found."; \
	fi

db-rollback-one: # support rollback the last migration version
	goose -dir db/migrations $$DATABASE_DRIVER "$$DATABASE_DRIVER://$$DATABASE_USERNAME:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_NAME?sslmode=disable" down

db-migrate: ## Perform all migration operations, commit file migrate & list migration status
	make migrate && make remove-migration-temp && make sqlc-migrations-commit && make db-migration-status

db-migration-status: ## Check the status of the migration, whether it is pending or has been migrated
	goose -dir db/migrations $$DATABASE_DRIVER "$$DATABASE_DRIVER://$$DATABASE_USERNAME:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_NAME?sslmode=disable" status

db-rollback: ## Rollback the number of db versions according to the parameter passed in "Ex: make db-rollback ROLLBACK_COUNT=2" -- the default is 1 version
	@if [ -z "$(ROLLBACK_COUNT)" ]; then \
		ROLLBACK_COUNT=1; \
	fi; \
	echo "Rolling back $(ROLLBACK_COUNT) migration(s)..."; \
	for i in $$(seq 1 $$ROLLBACK_COUNT); do \
		echo "Rollback round $$i..."; \
		make db-rollback-version; \
	done
	make db-migration-status

# sqlc migration commands
sqlc-db-columns: # Read of columns of database by table name
	@PGPASSWORD=$(DATABASE_PASSWORD) psql -U $(DATABASE_USERNAME) -d $(DATABASE_NAME) -h $(DATABASE_HOST) -p $(DATABASE_PORT) -t -A -c \
		"SELECT string_agg(column_name, ',') FROM information_schema.columns \
		WHERE table_name = '$(TABLE_NAME)' \
		AND column_name NOT IN ( \
			SELECT a.attname FROM pg_index i \
			JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey) \
			WHERE i.indrelid = '$(TABLE_NAME)'::regclass AND i.indisprimary \
		);"

sqlc-db-table-primary-key: # Read of primary key of database by table name
	@PGPASSWORD=$(DATABASE_PASSWORD) psql -U $(DATABASE_USERNAME) -d $(DATABASE_NAME) -h $(DATABASE_HOST) -p $(DATABASE_PORT) -t -A -c \
		"SELECT a.attname FROM pg_index i \
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey) \
		WHERE i.indrelid = '$(TABLE_NAME)'::regclass AND i.indisprimary;"

sqlc-generate-insert-query: ## Generate insert query SQL. Usage: make sqlc-generate-insert-query TABLE_NAME=users [FILENAME=insert_users]
	@if [ -z "$(TABLE_NAME)" ]; then \
		echo "❌ Missing TABLE_NAME. Ex: make sqlc-generate-insert-query TABLE_NAME=users [FILENAME=insert_users]"; \
		exit 1; \
	fi; \
	\
	OUT_FILE="db/queries/$(if $(FILENAME),$(FILENAME),$(TABLE_NAME)).sql"; \
	echo "📄 Output file: $$OUT_FILE"; \
	\
	COLUMNS=$$(make sqlc-db-columns TABLE_NAME=$(TABLE_NAME) | tail -n 1); \
	bash ./scripts/sqlc-generate-insert-query.sh $(TABLE_NAME) "$$COLUMNS" "$$OUT_FILE"

sqlc-generate-update-query: ## Generate update query SQL. Usage: make sqlc-generate-update-query TABLE_NAME=users [FILENAME=update_users]
	@if [ -z "$(TABLE_NAME)" ]; then \
		echo "❌ Missing TABLE_NAME. Ex: make sqlc-generate-update-query TABLE_NAME=users [FILENAME=update_users]"; \
		exit 1; \
	fi; \
	\
	OUT_FILE="db/queries/$(if $(FILENAME),$(FILENAME),$(TABLE_NAME)).sql"; \
	echo "📄 Output file: $$OUT_FILE"; \
	\
	COLUMNS=$$(make sqlc-db-columns TABLE_NAME=$(TABLE_NAME) | tail -n 1); \
	PRIMARY_KEY=$$(make sqlc-db-table-primary-key TABLE_NAME=$(TABLE_NAME) | tail -n 1); \
	bash ./scripts/sqlc-generate-update-query.sh $(TABLE_NAME) "$$COLUMNS" "$$PRIMARY_KEY" "$$OUT_FILE"

sqlc-generate-delete-query: ## Generate delete query SQL. Usage: make sqlc-generate-delete-query TABLE_NAME=users [FILENAME=delete_users]
	@if [ -z "$(TABLE_NAME)" ]; then \
		echo "❌ Missing TABLE_NAME. Ex: make sqlc-generate-delete-query TABLE_NAME=users [FILENAME=delete_users]"; \
		exit 1; \
	fi; \
	\
	OUT_FILE="db/queries/$(if $(FILENAME),$(FILENAME),$(TABLE_NAME)).sql"; \
	echo "📄 Output file: $$OUT_FILE"; \
	\
	PRIMARY_KEY=$$(make sqlc-db-table-primary-key TABLE_NAME=$(TABLE_NAME) | tail -n 1); \
	bash ./scripts/sqlc-generate-delete-query.sh $(TABLE_NAME) "$$PRIMARY_KEY" "$$OUT_FILE"

# Go generate commands
go-generate-module: ## Generate module. Usage: make go-generate-module MODULE_NAME=users [PACKAGE_DIR=users]
	bash ./scripts/go-generate-module.sh $(MODULE_NAME) $(PACKAGE_DIR)

# help
help: 
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*## "; print ""} /^[a-zA-Z_-]+:.*## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
