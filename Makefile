.PHONY: install start-backend stop-backend start-mobile clean-mobile update-memecoins test-backend test-mobile

# Install all dependencies
install:
	@echo "Installing backend dependencies..."
	go mod tidy
	@echo "Installing mobile app dependencies..."
	cd MemeTraderMobile && npm install

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
	cd MemeTraderMobile && npx expo start

# Clean mobile app build
clean-mobile:
	@echo "Cleaning mobile app build..."
	cd MemeTraderMobile && rm -rf node_modules
	cd MemeTraderMobile && npm install

# Update meme coins
update-memecoins:
	@echo "Updating meme coins..."
	curl -X POST http://localhost:8080/api/v1/memecoins/update

# Run backend tests
test-backend:
	@echo "Running backend tests..."
	go test ./... -v

# Run mobile app tests
test-mobile:
	@echo "Running mobile app tests..."
	cd MemeTraderMobile && npm test 