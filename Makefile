build:
	go build -v ./cmd/app

run:
	go run ./cmd/app/main.go --config=./configs/config.json

dbrun:
	docker run --name pizzasoft-pg-13.3 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=test -d postgres:13.3