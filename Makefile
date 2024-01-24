postgres:
	docker run --name webapp -p 5432:5432  -e POSTGRES_USER="postgres"  -e POSTGRES_PASSWORD="postgres" -d postgres:latest

createdb:
	docker exec -it webapp createdb --username="postgres" --owner="postgres" "users"

migratecreate:
	migrate create -ext sql -dir pkg/db/migrations/ -seq init_schema

migrateup:
	 migrate -path pkg/db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/users?sslmode=disable" -verbose up

dropdb:
	docker exec -it postgres dropdb "users"


.PHONY: postgres createdb createtestdb dropdb migrateup migratecreate