# Go API with Fiber and Bun

## Dev environments

Start database server

```bash
docker volume create pg-todo-db-dev

docker run \
--name todo-db-dev \
-p 2345:5432 \
-e POSTGRES_PASSWORD=mysecretpassword \
-e POSTGRES_DB=todo-db \
-v pg-todo-db-dev:/var/lib/postgresql/data \
-d postgres:14-alpine
```

Create config.yaml

```yaml
port: 3000
dsn: 'postgres://postgres:mysecretpassword@localhost:2345/todo-db?sslmode=disable'
```

## Migration

### Linux & Mac

```bash
docker run --rm --network host -v $(pwd)/migrations:/migrations migrate/migrate -verbose -path=/migrations/ -database "postgresql://postgres:mysecretpassword@localhost:2345/todo-db?sslmode=disable" up
```

### Windows PowerShell

```bash
docker run --rm --network host -v $pwd/migrations:/migrations migrate/migrate -verbose -path=/migrations/ -database postgresql://postgres:mysecretpassword@localhost:2345/todo-db?sslmode=disable up
```
