BINARY_NAME=server
OUTPUT_DIR=bin
EXECUTABLE_PKG=./

.PHONY: build run test clean

build: 
	@echo "building..."
	@go build -o ./$(OUTPUT_DIR)/$(BINARY_NAME) ./$(EXECUTABLE_PKG)

run: build
	@echo "executing..."
	@./$(OUTPUT_DIR)/$(BINARY_NAME)