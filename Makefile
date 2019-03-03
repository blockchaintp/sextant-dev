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
	docker-compose exec frontend yarn run develop

.PHONY: api.cli
api.cli:
	docker-compose exec api bash

.PHONY: api.run
api.run:
	docker-compose exec api yarn run start

.PHONY: psql
psql:
	docker-compose exec postgres psql --user postgres

.PHONY: database.migrate
database.migrate:
	docker-compose exec api yarn run knex -- migrate:latest