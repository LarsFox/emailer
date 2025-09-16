package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/LarsFox/emailer/api"
	"github.com/LarsFox/emailer/mail"
	"github.com/LarsFox/emailer/proto"
	"github.com/LarsFox/emailer/tg"
	"github.com/rollbar/rollbar-go"
)

var (
	emailUsername string
	emailPassword string
	emailHost     string
	emailPort     int64
	rollbarEnv    string
	rollbarToken  string
	tgBotToken    string
	serverPort    int64
)

// nolint:mnd
func flagParse() {
	flag.StringVar(&emailUsername, "email-username", "", "Account username")
	flag.StringVar(&emailPassword, "email-password", "", "Account password")
	flag.StringVar(&emailHost, "email-host", "", "Email server host")
	flag.Int64Var(&emailPort, "email-port", 587, "Email server port")
	flag.StringVar(&rollbarEnv, "rollbar-env", "dev", "Rollbar enb")
	flag.StringVar(&rollbarToken, "rollbar-token", "", "Rollbar token")
	flag.StringVar(&tgBotToken, "tg-bot-token", "", "TG Bot token")
	flag.Int64Var(&serverPort, "serve-port", 8080, "Application serving port")
	flag.Parse()
}

func main() {
	flagParse()
	port := fmt.Sprintf(":%d", serverPort)

	rollbar.SetCodeVersion(CommitHash)
	rollbar.SetEnvironment(rollbarEnv)
	rollbar.SetServerHost("emailer")
	rollbar.SetToken(rollbarToken)

	lc := &net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	emailClient := mail.NewClient(emailUsername, emailPassword, emailHost, emailPort)
	tgClient := tg.NewClient(tgBotToken)
	serv := grpc.NewServer()
	proto.RegisterEmailerServer(serv, api.NewServer(emailClient, tgClient))

	log.Printf("Listening %s...", port)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CommitHash — commit hash, used with -ldflags.
var CommitHash string
