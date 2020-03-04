package api

import (
	"context"
	"log"

	"github.com/LarsFox/emailer/proto"
)

type mailer interface {
	SendOneEmail(from, fromName, to, subj, msg string) error
}

// Server serves requests over gRPC.
type Server struct {
	mailer mailer
}

// NewServer returns new server.
func NewServer(mailer mailer) *Server {
	return &Server{mailer: mailer}
}

// SendOneEmail sends a single email.
func (s *Server) SendOneEmail(_ context.Context, in *proto.SendOneEmailRequest) (*proto.SendOneEmailResponse, error) {
	if err := s.mailer.SendOneEmail(in.From, in.FromName, in.To, in.Subject, in.Text); err != nil {
		log.Println(err)
		return &proto.SendOneEmailResponse{ErrorCode: 1}, nil
	}
	return &proto.SendOneEmailResponse{}, nil
}
