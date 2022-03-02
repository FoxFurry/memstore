build:
	go build -o bin/gostore

run:
	go run main.go

test:
	go test ./...

bench:
	go test -bench=. ./tests/