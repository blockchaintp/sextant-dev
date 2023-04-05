.PHONY: build
build:
	docker-compose build

.PHONY: dev
dev:
	docker-compose up -d

.PHONY: frontend.cli
frontend.cli:
	docker-compose exec frontend bash

.PHONY: frontend.run
frontend.run:
	docker-compose exec frontend npm run develop

.PHONY: api.cli
api.cli:
	docker-compose exec api bash

.PHONY: api.run
api.run:
	docker-compose exec api npm run serve

.PHONY: psql
psql:
	docker-compose exec postgres psql --user postgres

.PHONY: database.migrate
database.migrate:
	docker-compose exec api npm run knex -- migrate:latest

.PHONY: rebuild.api
rebuild.api:
	docker-compose down
	docker rmi sextant-dev-api || true
	docker-compose up -d
	docker-compose exec api bash
