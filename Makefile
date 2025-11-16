.PHONY: start_integration stop_integration integration


REPO_NAME := avito_test_task

RED := \033[0;31m
GREEN := \033[0;32m
NC := \033[0m

start_integration:
	@echo "starting integration tests services"
	@mkdir -p test/integration/coverage && \
	chmod 777 test/integration/coverage && \
	docker compose down && \
	docker compose pull && \
	docker compose up -d

stop_integration:
	@echo "stopping integration tests services"
	@cd tests/integration/docker && docker compose down

integration:
	@echo "run integration tests"
	@$(MAKE) start_integration
	@touch ./tests/integration/integration.lock
	@if ! go test ./tests/integration/... -v; then \
		rm -f ./tests/integration/integration.lock; \
		$(MAKE) stop_integration; \
		echo "$(RED)integration tests failed$(NC)"; \
		exit 1; \
	else \
		rm -f ./tests/integration/integration.lock; \
		$(MAKE) stop_integration; \
		echo "$(GREEN)integration tests passed$(NC)"; \
	fi
