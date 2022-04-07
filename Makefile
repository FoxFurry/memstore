run:
	go run main.go serve -p=8080

build:
	go build -o bin/gostore

test:
	go test ./...

bench:
	go test -bench=. ./tests/
