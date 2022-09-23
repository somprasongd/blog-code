DSN=postgres://postgres:mysecretpassword@localhost:2345/todo-db?sslmode=disable

export DSN

dev-up:
	docker volume create pg-todo-db-dev \
	&& docker run \
	--name todo-db-dev \
	-p 2345:5432 \
	-e POSTGRES_PASSWORD=mysecretpassword \
	-e POSTGRES_DB=todo-db \
	-v pg-todo-db:/var/lib/postgresql/data \
	-d postgres:14-alpine

dev-down:
	docker stop todo-db-dev && docker rm todo-db-dev

migrate:
	docker run --rm --network host -v $(pwd)/migrations:/migrations migrate/migrate -verbose -path=/migrations/ -database postgres://postgres:mysecretpassword@localhost:2345/todo-db?sslmode=disable up

migrate-win:
	docker run --rm --network host -v $pwd/migrations:/migrations migrate/migrate -verbose -path=/migrations/ -database postgres://postgres:mysecretpassword@localhost:2345/todo-db?sslmode=disable up

run:
	go run cmd/main.go