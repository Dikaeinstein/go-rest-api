BINARY_NAME=go-rest-api

## Build binary
build:
	go build

## Build and install binary
install:
	go build -install

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
