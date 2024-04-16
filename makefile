build:
	go build -o bin/main main.go

run: build
	go run main.go