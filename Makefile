postgres: 
	docker run --name brm-postgres-con -p 5433:5432 -e POSTGRES_USER=dev_user@brm -e POSTGRES_PASSWORD=pass1234 -d postgres:14-alpine

createdb:
	docker exec -it brm-postgres-con createdb --username=dev_user@brm --owner=dev_user@brm brm_test

dropdb:
	docker exec -it brm-postgres-con dropdb brm_test

migrateup:
	migrate -path db/migration -database "postgresql://dev_user@brm:pass1234@localhost:5433/brm_test?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://dev_user@brm:pass1234@localhost:5433/brm_test?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://dev_user@brm:pass1234@localhost:5433/brm_test?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://dev_user@brm:pass1234@localhost:5433/brm_test?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

go-test:
	go test -v -cover ./...

server: 
	go run main.go -env app.env

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc go-test server