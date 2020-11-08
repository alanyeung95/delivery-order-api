.DEFAULT_GOAL := run

NETWORK := network

.PHONY: network
network:
	docker network create $(NETWORK) || true

.PHONY: run
run: network
	docker-compose up

.PHONY: unit-test
unit-test: network
	docker-compose -f docker-compose.unit.test.yml up

.PHONY: integration-test
integration-test: network
	docker-compose -f docker-compose.yml -f docker-compose.integration.test.yml up  --abort-on-container-exit

.PHONY: clean
clean: 
	docker-compose down --remove-orphans