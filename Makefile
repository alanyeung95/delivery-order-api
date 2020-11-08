.PHONY: unit-test
unit-test:
	docker-compose -f docker-compose.unit.test.yml up

.PHONY: integration-test
integration-test:
	docker-compose -f docker-compose.yml -f docker-compose.integration.test.yml up  --abort-on-container-exit

.PHONY: clean
clean: 
	docker-compose down --remove-orphans