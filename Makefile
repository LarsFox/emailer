include .env

default: generate lint run

VERSION := "-X main.CommitHash=$(shell git rev-parse --short HEAD)_$(shell date -u +%Y-%m-%d_%H:%M:%S)"

run:
	@go run -ldflags $(VERSION) -race ./cmd/main.go \
		--email-username=$(EMAILER_USERNAME) \
		--email-password=$(EMAILER_PASSWORD) \
		--email-host=$(EMAILER_SERVER_HOST) \
		--rollbar-env=$(ROLLBAR_ENV) \
		--rollbar-token=$(ROLLBAR_TOKEN) \
		--tg-bot-token=$(TG_BOT_TOKEN) \
		--serve-port=$(EMAILER_PORT)

lint:
	@golangci-lint run

generate:
	@go generate ./...

protoc:
	@cd proto && protoc --go_out=plugins=grpc:. *.proto

b:
	go build -o emailer -ldflags $(VERSION) ./cmd/main.go

bp:
	env GOOS=linux GOARCH=386 go build -o emailer -ldflags $(VERSION) ./cmd/main.go

up:
	scp emailer $(REMOTE_HOST_PATH)/emailer_test
