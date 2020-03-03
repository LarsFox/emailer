package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/LarsFox/emailer/api"
	"github.com/LarsFox/emailer/mail"
	"github.com/LarsFox/emailer/proto"
)

var (
	emailUsername string
	emailPassword string
	emailHost     string
	emailPort     int64
	serverPort    int64
)

func init() {
	flag.StringVar(&emailUsername, "email-username", "", "Account username")
	flag.StringVar(&emailPassword, "email-password", "", "Account password")
	flag.StringVar(&emailHost, "email-host", "", "Email server host")
	flag.Int64Var(&emailPort, "email-port", 587, "Email server port")
	flag.Int64Var(&serverPort, "serve-port", 8080, "Application serving port")
	flag.Parse()
}

func main() {
	port := fmt.Sprintf(":%d", serverPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	emailClient := mail.NewClient(emailUsername, emailPassword, emailHost, emailPort)
	serv := grpc.NewServer()
	proto.RegisterEmailerServer(serv, api.NewServer(emailClient))

	log.Printf("Listening %s...", port)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
