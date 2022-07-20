DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
PWD=@2Z!I95nj%rx2

network:
	docker network create bank-network
postgres:
	docker run --name postgres_tutorial --network bank-network  -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres_tutorial createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres_tutorial dropdb simple_bank

migrateup:
	migrate  -path db/migration -database "postgresql://root:2ZI95njrx2@database-1.c6d1dadhgyqb.us-east-1.rds.amazonaws.com:5432/simplebank" -verbose up

migrateup1:
	migrate  -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate  -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate  -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgress -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/parin/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1  migratedown1 sqlc test server mock db_docs db_schema