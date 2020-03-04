include .env

default: generate lint run

run:
	@go run -race ./cmd/main.go \
		--email-username=$(EMAILER_USERNAME) \
		--email-password=$(EMAILER_PASSWORD) \
		--email-host=$(EMAILER_SERVER_HOST) \
		--serve-port=$(EMAILER_PORT)

lint:
	@golangci-lint run

generate:
	@go generate ./...

protoc:
	@cd proto && protoc --go_out=plugins=grpc:. *.proto

build:
	@mkdir -p .tmp
	@go build -o .tmp/emailer -race ./cmd/main.go
