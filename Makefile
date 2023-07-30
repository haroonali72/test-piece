# Build the Docker image and start the application
start:
	docker-compose up --build -d

# Stop and remove the Docker containers
stop:
	docker-compose down

# Show logs from the application container
logs:
	docker-compose logs -f app
