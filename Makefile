APP_NAME := tjournal.exe

tidy:
	@echo "Tidying up..."
	@go mod tidy
	@go mod vendor

build:
	@echo "building..."
	@go build -o ./bin/${APP_NAME} ./src/main.go

run: tidy build
	@./bin/${APP_NAME}

dev: build
	@./bin/${APP_NAME}

