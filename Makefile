run:
	go run cmd/main.go
migrateup:
	migrate -path db/migration -database "postgres://postgres:Lovego@localhost:5432/test_aprel?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:Lovego@localhost:5432/test_aprel?sslmode=disable" -verbose down
migrateforce:
	migrate -path db/migration -database "postgres://postgres:Lovego@localhost:5432/test_aprel?sslmode=disable" force 1
	make migratedown
.phony: migrateup migratedown
