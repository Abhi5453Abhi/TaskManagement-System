.PHONY: backend-test frontend-test ci backend-build frontend-build backend-run frontend-run clean

# Backend commands
backend-test:
	cd backend && go test ./... -v

backend-build:
	cd backend && go build -o bin/api cmd/api/main.go

backend-run:
	cd backend && go run cmd/api/main.go

# Frontend commands
frontend-test:
	cd frontend && npm test

frontend-build:
	cd frontend && npm run build

frontend-run:
	cd frontend && npm run dev

# Combined commands
ci: backend-test frontend-test
	@echo "All tests passed!"

# Development commands
dev: backend-run frontend-run

# Clean up
clean:
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -f backend/tasks.db

# Install dependencies
install:
	cd backend && go mod tidy
	cd frontend && npm install

# Check line count
check-loc:
	@echo "Checking line count..."
	@find . -name "*.go" -o -name "*.ts" -o -name "*.tsx" | grep -v node_modules | grep -v vendor | xargs wc -l | tail -1
