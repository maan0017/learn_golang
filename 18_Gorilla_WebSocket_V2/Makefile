EXECUTABLE_NAME = server
OUTPUT_DIR = bin
MAIN_PKG = ./

.PHONY: build run

build:
	@echo "building 18_Gorilla_WebSocket_V2 ..."
	@go build -o ./$(OUTPUT_DIR)/$(EXECUTABLE_NAME) ./$(MAIN_PKG)

run: build
	@echo "executing...."
	@./$(OUTPUT_DIR)/$(EXECUTABLE_NAME)