dev.up:
	@echo "Starting local environment..."
	docker compose up --remove-orphans

dev.down:
	@echo "Stopping local environment..."
	docker compose down --volumes

test.mock:
	@echo "Generating mock files..."
	mockgen -source=internal/service/service.go -destination=mocks/mock_service.go -package=mocks