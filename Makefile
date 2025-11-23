.PHONY: build run test test-e2e lint clean docker-up docker-down

APP_NAME = app

test:
	go test -v ./...

test-e2e:
	@echo "Running E2E tests..."
	go test -v -tags=e2e -count=1 ./tests/e2e_test/...

lint:
	@echo "[LINT] Checking for golangci-lint..."
	@which golangci-lint >/dev/null 2>&1 || { \
		echo "[INFO] Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
	@echo "[LINT] Running linter..."
	golangci-lint run --fix ./...
	
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

