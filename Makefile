postgres:
	docker run --name simple_pg -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine3.21

createdb:
	docker exec -it simple_pg createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple_pg dropdb simple_bank

migrate_up: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc test