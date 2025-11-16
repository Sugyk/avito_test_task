.PHONY: start_integration stop_integration integration


REPO_NAME := avito_test_task

RED := \033[0;31m
GREEN := \033[0;32m
NC := \033[0m

start_integration:
	@echo "starting integration tests services"
	@docker compose down && \
	docker compose build && \
	docker compose pull && \
	docker compose up -d

stop_integration:
	@echo "stopping integration tests services"
	@docker compose down

integration:
	@echo "run integration tests"
	@$(MAKE) start_integration
	@if ! go test ./tests/integration/... -v; then \
		$(MAKE) stop_integration; \
		echo "$(RED)integration tests failed$(NC)"; \
		exit 1; \
	else \
		$(MAKE) stop_integration; \
		echo "$(GREEN)integration tests passed$(NC)"; \
	fi

unit:
	@echo "run unit tests"
	@go test ./internal/...
