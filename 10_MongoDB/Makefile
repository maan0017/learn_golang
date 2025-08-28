# Makefile for Go project

BINARY_NAME=server
OUTPUT_DIR=bin
MAIN_PKG=./cmd/main

.PHONY: build run test clean

build:
	@echo "Building..."
	@go build -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_PKG)

run: build
	@echo "Running..."
	@./$(OUTPUT_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)
