package api

import (
	"context"
	"errors"
	"log"

	"github.com/rollbar/rollbar-go"

	"github.com/LarsFox/emailer/entities"
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
	if err := s.mailer.SendOneEmail(in.GetFrom(), in.GetFromName(), in.GetTo(), in.GetSubject(), in.GetText()); err != nil {
		log.Println(err)
		return &proto.SendOneEmailResponse{ErrorCode: 1}, nil
	}
	return &proto.SendOneEmailResponse{}, nil
}

// SendOneTGMessage sends a single TGMessage.
func (s *Server) SendOneTGMessage(ctx context.Context, in *proto.SendOneTGMessageRequest) (*proto.SendOneTGMessageResponse, error) {
	if err := s.telegramer.SendMessage(ctx, in.GetTo(), in.GetText()); err != nil {
		var tgErr *entities.TGError
		if errors.As(err, &tgErr) {
			return &proto.SendOneTGMessageResponse{ErrorCode: tgErr.Code}, nil
		}
		rollbar.Error(err)
		return &proto.SendOneTGMessageResponse{ErrorCode: 1}, nil
	}
	return &proto.SendOneTGMessageResponse{}, nil
}
