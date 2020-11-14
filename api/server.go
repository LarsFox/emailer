package api

import (
	"context"
	"log"

	"github.com/bugsnag/bugsnag-go"

	"github.com/LarsFox/emailer/common"
	"github.com/LarsFox/emailer/proto"
)

type mailer interface {
	SendOneEmail(from, fromName, to, subj, msg string) error
}

type telegramer interface {
	SendMessage(ctx context.Context, to int64, text string) error
}

// Server serves requests over gRPC.
type Server struct {
	mailer     mailer
	telegramer telegramer
}

// NewServer returns new server.
func NewServer(mailer mailer, telegramer telegramer) *Server {
	return &Server{mailer: mailer, telegramer: telegramer}
}

// SendOneEmail sends a single email.
func (s *Server) SendOneEmail(_ context.Context, in *proto.SendOneEmailRequest) (*proto.SendOneEmailResponse, error) {
	if err := s.mailer.SendOneEmail(in.From, in.FromName, in.To, in.Subject, in.Text); err != nil {
		log.Println(err)
		return &proto.SendOneEmailResponse{ErrorCode: 1}, nil
	}
	return &proto.SendOneEmailResponse{}, nil
}

// SendOneTGMessage sends a single TGMessage.
func (s *Server) SendOneTGMessage(ctx context.Context, in *proto.SendOneTGMessageRequest) (*proto.SendOneTGMessageResponse, error) {
	if err := s.telegramer.SendMessage(ctx, in.To, in.Text); err != nil {
		if e, ok := err.(*common.TGError); ok {
			return &proto.SendOneTGMessageResponse{ErrorCode: e.Code}, nil
		}
		bugsnag.Notify(err)
		return &proto.SendOneTGMessageResponse{ErrorCode: 1}, nil
	}
	return &proto.SendOneTGMessageResponse{}, nil
}
