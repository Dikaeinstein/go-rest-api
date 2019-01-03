BINARY_NAME=go-rest-api

## Fetch dependencies
install:
	go get -t -v ./...

## Run tests
test:
	APP_ENV=test go test -v ./...

## Build binary
build:
	go build

## Execute binary
run:build
	./$(BINARY_NAME)

.PHONY: clean
## Remove binary
clean:
	rm -f $(BINARY_NAME)
