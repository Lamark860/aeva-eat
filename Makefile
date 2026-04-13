.PHONY: lint lint-go lint-vue test test-go build up down logs

# === Lint ===
lint: lint-go lint-vue

lint-go:
	@echo "==> Go vet..."
	cd backend && go vet ./...
	@echo "==> Go build check..."
	cd backend && go build ./...

lint-vue:
	@echo "==> Vue ESLint..."
	docker exec aeva-frontend npx eslint src --ext .vue,.js --max-warnings 0
	@echo "==> Vue build check..."
	docker exec aeva-frontend npx vite build --mode development 2>&1 | tail -5

# === Test ===
test: test-go

test-go:
	@echo "==> Go tests..."
	cd backend && go test ./... -count=1 -v

# === Build ===
build:
	docker compose build

build-no-cache:
	docker compose build --no-cache

# === Docker ===
up:
	docker compose up -d

up-build:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f --tail=50

logs-backend:
	docker compose logs -f --tail=50 backend

# === All checks ===
check: lint test
	@echo "==> All checks passed!"
