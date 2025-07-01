ATLAS = atlas

MIGRATION_NAME ?= new_migration
ENV ?= local

postgres:
	docker run --name simpleblank_db -e POSTGRES_PASSWORD=030603 -e POSTGRES_USER=admin -e POSTGRES_DB=core_bank -p 5000:5432 -d postgres

createdb:
	docker exec -it simpleblank_db createdb --username=admin --owner=admin core_bank

dropdb:
	docker exec -it simpleblank_db dropdb --username=admin core_bank

migrate-diff:
	@echo "-> Generando migración: $(MIGRATION_NAME)"
	$(ATLAS) migrate diff $(MIGRATION_NAME) --env $(ENV)

migrate-apply:
	@echo "-> Aplicando migraciones..."
	$(ATLAS) migrate apply --env $(ENV)

migrate-status:
	@echo "-> Verificando estado de las migraciones..."
	$(ATLAS) migrate status --env $(ENV)

migrate-new:
	@echo "📁 Creando migración vacía: $(MIGRATION_NAME)"
	$(ATLAS) migrate new $(MIGRATION_NAME) --env $(ENV)

sqlc:
	@echo "-> Generando código SQL..."
	sqlc generate

migrate-hash:
	@echo "-> Recalculando hashes de migraciones (atlas.sum)..."
	atlas migrate hash --env $(ENV)

test:
	@echo "-> Ejecutando pruebas..."
	go test -v -cover ./...

.PHONY: migrate-hash sqlc