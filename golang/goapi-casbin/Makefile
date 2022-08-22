SERVICE_NAME=Auth-Api
SERVICE_IMAGE=somprasongd/auth-api
SERVICE_VERSION=1.0.0

export SERVICE_NAME
export SERVICE_IMAGE
export SERVICE_VERSION

doc:
	@echo "---Generate doc files $(SERVICE_NAME)---"
	swag init -g pkg/module/module.go -o pkg/docs

dev-up:
	@echo "---Start Dev $(SERVICE_NAME) Environtment---"
	docker-compose -p auth-api-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml up -d

dev-down:
	@echo "---Stop Dev $(SERVICE_NAME) Environtment---"
	docker-compose -p auth-api-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml down

dev-down-v:
	@echo "---Stop Dev $(SERVICE_NAME) Environtment---"
	docker-compose -p auth-api-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml down -v

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

prod-down-v:
	@echo "---Stop Prod $(SERVICE_NAME)---"
	docker-compose -p task-api-prod -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.prod.yml down -v

test-handler:
	@echo "---Start test auth handler---"
	go test -v goapi/pkg/module/auth/handler

test-service:
	@echo "---Start test auth service---"
	go test -v goapi/pkg/module/auth/core/service

test-all:
	@echo "---Start test all module---"
	go test -v goapi/pkg/module/...

migrate-create:
	@echo "---Creating migration files---"
	go run cmd/migrate/main.go create $(NAME) sql

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-force:
	go run cmd/migrate/main.go force $(VERSION)

migrate-status:
	go run cmd/migrate/main.go status

migrate-version:
	go run cmd/migrate/main.go version