.PHONY: install start-backend stop-backend start-mobile clean-mobile update-memecoins test-backend test-mobile test-blockchain test-providers test-coverage start-services test-integration test-integration-short start-mobile-ios start-mobile-android install-mobile-deps

# Install all dependencies
install:
	@echo "Installing backend dependencies..."
	go mod tidy
	@echo "Installing mobile app dependencies..."
	@make install-mobile-deps

# Install mobile dependencies
install-mobile-deps:
	@echo "Installing mobile app dependencies..."
	cd MemeTraderMobileNew && npm install --legacy-peer-deps
	cd MemeTraderMobileNew && npx expo install

# Start all services
start-services:
	@echo "Starting all services..."
	@make start-backend
	@echo "Waiting for backend services to be ready..."
	@sleep 5
	@make start-mobile

# Start backend services
start-backend:
	@echo "Starting backend services..."
	docker-compose down -v
	docker-compose up --build -d

# Stop backend services
stop-backend:
	@echo "Stopping backend services..."
	docker-compose down -v

# Start mobile app
start-mobile:
	@echo "Starting mobile app..."
	cd MemeTraderMobileNew && npx expo start

# Start mobile app on iOS simulator
start-mobile-ios:
	@echo "Starting mobile app on iOS simulator..."
	cd MemeTraderMobileNew && npx expo start --ios

# Start mobile app on Android emulator
start-mobile-android:
	@echo "Starting mobile app on Android emulator..."
	cd MemeTraderMobileNew && npx expo start --android

# Clean mobile app build
clean-mobile:
	@echo "Cleaning mobile app build..."
	cd MemeTraderMobileNew && rm -rf node_modules
	cd MemeTraderMobileNew && rm -rf .expo
	cd MemeTraderMobileNew && npm cache clean --force
	@make install-mobile-deps

# Update meme coins
update-memecoins:
	@echo "Updating meme coins..."
	curl -X POST http://localhost:8080/api/v1/memecoins/update

# Run all backend tests
test-backend:
	@echo "Running backend tests..."
	go clean -testcache
	go test ./... -v

# Run blockchain provider tests
test-blockchain:
	@echo "Running blockchain provider tests..."
	go clean -testcache
	go test ./internal/blockchain/... -v

# Run specific provider tests
test-providers:
	@echo "Running provider tests..."
	go clean -testcache
	go test ./internal/blockchain/solana/... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go clean -testcache
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run mobile app tests
test-mobile:
	@echo "Running mobile app tests..."
	cd MemeTraderMobileNew && npm test

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go clean -testcache
	go test ./internal/blockchain/providers/... -v

# Run integration tests in short mode (skips long-running tests)
test-integration-short:
	@echo "Running integration tests in short mode..."
	go clean -testcache
	go test -short ./internal/blockchain/providers/... -v 