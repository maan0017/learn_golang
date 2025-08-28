EXECUTABLE_NAME=server.exe
OUTPUT_DIR=bin
MAIN_PKG=./

.PHONY: build run clean

build: 
	@echo "building..."
	@go build -o ./$(OUTPUT_DIR)/$(EXECUTABLE_NAME) ./$(MAIN_PKG)

run: build
	@echo "executing..."
	@./$(OUTPUT_DIR)/$(EXECUTABLE_NAME)

clean: 
	@echo "cleaning..."
	@rm -rf ./$(OUTPUT_DIR)