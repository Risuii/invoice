build:
	go build cmd/main.go

run:
	go mod tidy
	go run cmd/main.go

migrate.build:
	migrate create -ext sql -dir migration/sql/ -seq init_mg

migrate.up:
	go run migration/main/main.go up

migrate.rollback:
	go run migration/main/main.go rollback

test:
	go test -coverprofile cover.out ./src/...
	go tool cover -html=cover.out
	
	

