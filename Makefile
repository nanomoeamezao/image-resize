BINARY_NAME=bin/resizer

.PHONY: build 
build:
	go build -o $(BINARY_NAME) main.go
