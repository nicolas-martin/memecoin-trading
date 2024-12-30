.PHONY: help install start-backend stop-backend start-mobile install-watchman clean-mobile update-memecoins test-backend test-mobile

# Default target when just running 'make'
help:
	@echo "Available commands:"
	@echo "  make install           - Install all dependencies (backend & mobile)"
	@echo "  make start-backend    - Start the backend services with Docker"
	@echo "  make stop-backend     - Stop the backend services"
	@echo "  make start-mobile     - Start the Expo development server"
	@echo "  make install-watchman - Install Watchman (required for mobile development)"
	@echo "  make clean-mobile     - Clean mobile app build and dependencies"
	@echo "  make update-memecoins - Trigger memecoins update from API"
	@echo "  make test-backend     - Run backend tests"
	@echo "  make test-mobile      - Run mobile app tests"

# Installation commands
install: install-watchman
	@echo "Installing backend dependencies..."
	go mod tidy
	@echo "Installing mobile app dependencies..."
	cd MemeTraderMobileNew && npm install --legacy-peer-deps

install-watchman:
	@echo "Installing Watchman..."
	brew install watchman

# Backend commands
start-backend:
	@echo "Starting backend services..."
	docker-compose down -v
	docker-compose up --build -d

stop-backend:
	@echo "Stopping backend services..."
	docker-compose down

update-memecoins:
	@echo "Updating memecoins..."
	curl -X POST http://localhost:8080/api/v1/memecoins/update

test-backend:
	@echo "Running backend tests..."
	go test ./... -v

# Mobile commands
start-mobile:
	@echo "Starting Expo development server..."
	cd MemeTraderMobileNew && npx expo start

clean-mobile:
	@echo "Cleaning mobile app..."
	cd MemeTraderMobileNew && rm -rf node_modules
	cd MemeTraderMobileNew && rm -rf .expo
	cd MemeTraderMobileNew && npm cache clean --force

test-mobile:
	@echo "Running mobile app tests..."
	cd MemeTraderMobileNew && npm test 