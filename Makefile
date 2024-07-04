DB_SOURCE=postgresql://root:secret@localhost:5432/timetracker?sslmode=disable
USER=root
PASSWORD=secret



.PHONY: migrateup migratedown migrateup1 migratedown1 start

createdb:
	docker exec -it ttcontainer createdb --username=$(USER) --owner=$(USER) timetracker

dropdb:
	docker exec -it ttcontainer dropdb timetracker

postgres:
	docker run --name ttcontainer -p 5432:5432 -e POSTGRES_USER=$(USER) -e POSTGRES_PASSWORD=$(PASSWORD) -d postgres:15-alpine

start:
	docker start ttcontainer

migrateup:
	migrate -path ./migrations -database "${DB_SOURCE}" -verbose up

migrateup1:
	migrate -path ./migrations -database "${DB_SOURCE}" -verbose up 1

migratedown:
	migrate -path ./migrations -database "${DB_SOURCE}" -verbose down

migratedown1:
	migrate -path ./migrations -database "${DB_SOURCE}" -verbose down 1

start:
	go run cmd/main.go