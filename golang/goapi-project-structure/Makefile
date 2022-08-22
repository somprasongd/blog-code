SERVICE_NAME=Todo-Api
SERVICE_IMAGE=somprasongd/todo-api
SERVICE_VERSION=1.0.0

export SERVICE_NAME
export SERVICE_IMAGE
export SERVICE_VERSION

dev-up:
	@echo "---Start Dev $(SERVICE_NAME) Environtment---"
	docker-compose -p todo-api-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml up -d

dev-down:
	@echo "---Stop Dev $(SERVICE_NAME) Environtment---"
	docker-compose -p todo-api-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml down

dev:
	@echo "---Start Dev $(SERVICE_NAME)---"
	go run cmd/api/main.go

d-build:
	@echo "---Build $(SERVICE_NAME) $(SERVICE_IMAGE):$(SERVICE_VERSION)---"
	docker build -t $(SERVICE_IMAGE):$(SERVICE_VERSION) -f deploy/Dockerfile .

d-build-debug:
	@echo "---Build $(SERVICE_NAME) $(SERVICE_IMAGE):$(SERVICE_VERSION)---"
	docker build --progress plain -t $(SERVICE_IMAGE):$(SERVICE_VERSION) -f deploy/Dockerfile .

prod-up:
	@echo "---Start Prod $(SERVICE_NAME)---"
	docker-compose -p task-api-prod -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.prod.yml up -d

prod-down:
	@echo "---Stop Prod $(SERVICE_NAME)---"
	docker-compose -p task-api-prod -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.prod.yml down