build:
	mkdir -p bin
	go build -o bin/art main.go

test:
	go test ./...