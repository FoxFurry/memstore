run:
	go run main.go

build:
	go build -o bin/gostore

test:
	go test ./...

bench:
	go test -bench=. ./tests/
