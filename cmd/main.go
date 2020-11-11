package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/bugsnag/bugsnag-go"
	"google.golang.org/grpc"

	"github.com/LarsFox/emailer/api"
	"github.com/LarsFox/emailer/mail"
	"github.com/LarsFox/emailer/proto"
	"github.com/LarsFox/emailer/tg"
)

var (
	bugsnagKey    string
	emailUsername string
	emailPassword string
	emailHost     string
	emailPort     int64
	tgBotToken    string
	serverPort    int64
)

func flagParse() {
	flag.StringVar(&bugsnagKey, "bugsnag-key", "", "Bugsnag key")
	flag.StringVar(&emailUsername, "email-username", "", "Account username")
	flag.StringVar(&emailPassword, "email-password", "", "Account password")
	flag.StringVar(&emailHost, "email-host", "", "Email server host")
	flag.Int64Var(&emailPort, "email-port", 587, "Email server port")
	flag.StringVar(&tgBotToken, "tg-bot-token", "", "TG Bot token")
	flag.Int64Var(&serverPort, "serve-port", 8080, "Application serving port")
	flag.Parse()
}

func main() {
	flagParse()
	port := fmt.Sprintf(":%d", serverPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer bugsnag.Recover()
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       bugsnagKey,
		AppVersion:   getAppVersion(),
		PanicHandler: func() {},
	})

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

func getAppVersion() string {
	semver := "1.0.0"
	if CommitHash == "" {
		return semver
	}
	return fmt.Sprintf("%s-%s", semver, CommitHash)
}
