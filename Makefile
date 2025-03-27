postgres:
	docker run --name postgresCont -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
dbStart:
	docker start postgresCont
dbStop:
	docker stop postgresCont
dbRemove:
	docker rm postgresCont
createdb:
	docker exec -it postgresCont createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgresCont dropdb simple_bank	
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock :
	mockgen -package mockdb -destination db/mock/store.go github.com/Sidsha242/simple_bank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server



#truncate:
    #docker exec -i postgresCont psql -U root -d simple_bank -f /db/truncate_all_tables.sql