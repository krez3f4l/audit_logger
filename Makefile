.SILENT:

build:
	source ./.env && go build -o app ./cmd/main.go

run: build
	./app

	