dev.up:
	@echo "Starting local environment..."
	docker compose up --remove-orphans

dev.down:
	@echo "Stopping local environment..."
	docker compose down --volumes

test:
	@echo "Run testing..."
	ginkgo -r -p --randomize-suites --randomize-all --fail-on-pending --trace --race --show-node-events -cover
