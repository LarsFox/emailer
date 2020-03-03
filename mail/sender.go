package mail

import (
	"fmt"
	"net/smtp"
)

// Client sends emails.
type Client struct {
	auth *smtp.Auth
	host string
	addr string
}

// NewClient returns a new client.
func NewClient(username, password, host string, port int64) *Client {
	auth := smtp.PlainAuth("", username, password, host)
	return &Client{auth: &auth, host: host, addr: fmt.Sprintf("%s:%d", host, port)}
}

var messageDummy = `To:%s
Subject: %s

%s`

// SendOneEmail sends one simple email.
func (c *Client) SendOneEmail(from, to, subj, msg string) error {
	text := fmt.Sprintf(messageDummy, to, subj, msg)
	return smtp.SendMail(c.addr, *c.auth, from, []string{to}, []byte(text))
}
