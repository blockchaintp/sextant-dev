.PHONY: build
build:
	docker-compose build

.PHONY: dev
dev:
	docker-compose up

.PHONY: frontend.cli
frontend.cli:
	docker-compose exec frontend bash
	
.PHONY: api.cli
api.cli:
	docker-compose exec api bash
