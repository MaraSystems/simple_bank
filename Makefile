postgres:
	docker run --name simple_pg -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine3.21

createdb:
	docker exec -it simple_pg createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple_pg dropdb simple_bank

migrate_create:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate_up: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable" -verbose up

migrate_up_latest: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable" -verbose up 1

migrate_down:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable" -verbose down

migrate_down_latest: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

serve:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/MaraSystems/simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrate_up migrate_down sqlc test mock