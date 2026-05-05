# Docker commands
db-up:
	docker-compose up -d

db-down:
	docker-compose down

db-logs:
	docker-compose logs -f db

# Backend commands
build-backend:
	cd backend && mkdir -p ./bin && go build -o ./bin/moneyKeeper ./cmd/...

run-backend:
	cd backend && go run ./cmd/main.go

# Frontend commands
run-frontend:
	cd frontend && npm run dev

# Full environment
dev: db-up
	@echo "Database is starting..."
	@echo "Starting backend and frontend..."
	# Using & to run in parallel, might want to use separate terminals though
	(cd backend && go run ./cmd/main.go) & (cd frontend && npm run dev)
