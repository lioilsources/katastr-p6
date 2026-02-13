.PHONY: run test build lint clean

run:
	go run cmd/server/main.go

test:
	go test -v ./...

build:
	go build -o bin/server cmd/server/main.go

lint:
	go vet ./...

clean:
	rm -rf bin/
