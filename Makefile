ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

POSTGRES_INTEGRATION := user=test password=test dbname=test host=localhost port=15432 sslmode=disable

MIGRATION_FOLDER=$(CURDIR)/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: run-test-environment
run-test-environment:
	docker compose up -d

.PHONY: stop-test-environment
stop-test-environment:
	docker compose down

.PHONY: run-unit-tests
run-unit-tests:
	go test ./...

.PHONY: run-integration-environment
run-integration-environment:
	docker compose -f docker-compose-integration.yaml up -d
	sleep 5
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_INTEGRATION)" up

.PHONY: stop-integration-environment
stop-integration-environment:
	docker compose -f docker-compose-integration.yaml down

.PHONY: run-integration-tests
run-integration-tests:
	DB_PORT=15432 go test -tags=integration ./...
