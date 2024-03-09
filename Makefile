# It runs the docker-compose up command in detached mode to start the services defined in the compose.yaml file
start-services:
	@echo "🌀Initializing services for take off..."
	@echo "♻️ Cleaning up previous containers..."
	@docker-compose -f compose.yaml up -d
	@echo "✅ Services are up and running! 🚀"

# It runs the docker-compose down command to stop and remove the services defined in the compose.yaml file.
stop-services:
	@echo "🛑 Stopping services..."
	@docker-compose -f compose.yaml down
	@echo "✅ Services are stopped! 🛑"
