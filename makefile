BINARY_NAME=go-rest-api

## Build binary
build:
	go build

## Fetch dependencies
install:
	go get -t -v ./...

## Run tests
test:
	APP_ENV=test go test -v ./...

## Execute binary
run:build
	./$(BINARY_NAME)

.PHONY: clean
## Remove binary
clean:
	rm -f $(BINARY_NAME)
