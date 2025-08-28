EXECUTABLE_NAME=server
OUTPUT_DIR=bin
MAIN_PKG=./cmd/server

.PHONY: build run backend frontend clean

## Build backend binary
build:
	@echo "building backend..."
	@cd ./backend && go build -o ./$(OUTPUT_DIR)/$(EXECUTABLE_NAME) $(MAIN_PKG)

## Run backend only
backend: build
	@echo "starting backend..."
	@cd ./backend && ./$(OUTPUT_DIR)/$(EXECUTABLE_NAME)

## Run frontend only
frontend:
	@echo "starting frontend..."
	@cd ./frontend && npm run dev

## Run both backend + frontend in parallel
run: build
	@echo "starting backend + frontend..."
	@cd ./backend && ./$(OUTPUT_DIR)/$(EXECUTABLE_NAME) & \
	cd ./frontend && npm run dev

## Clean build artifacts
clean:
	@echo "cleaning..."
	@cd ./backend && rm -rf $(OUTPUT_DIR)
