include .env

default: generate lint run

run:
	@go run -ldflags "-X main.CommitHash=`git rev-parse --short HEAD`" -race ./cmd/main.go \
		--bugsnag-key=$(BUGSNAG_KEY) \
		--email-username=$(EMAILER_USERNAME) \
		--email-password=$(EMAILER_PASSWORD) \
		--email-host=$(EMAILER_SERVER_HOST) \
		--tg-bot-token=$(TG_BOT_TOKEN) \
		--serve-port=$(EMAILER_PORT)

lint:
	@golangci-lint run

generate:
	@go generate ./...

protoc:
	@cd proto && protoc --go_out=plugins=grpc:. *.proto

b:
	go build -o emailer -ldflags "-X main.CommitHash=`git rev-parse HEAD`" -race ./cmd/main.go

bp:
	env GOOS=linux GOARCH=386 go build -o emailer -ldflags "-X main.CommitHash=`git rev-parse HEAD`" ./cmd/main.go

up:
	scp emailer $(MOTOVSKIKH_REMOTE_HOST_PATH)/emailer_test
