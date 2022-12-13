package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LarsFox/emailer/common"
)

// Client works with Telegram.
type Client struct {
	token string
}

type response struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description json.RawMessage `json:"description"` // only if not OK
	ErrorCode   int64           `json:"error_code"`
}

type sendMessage struct {
	ChatID            int64  `json:"chat_id"`
	DisableURLPreview bool   `json:"disable_web_page_preview"`
	Text              string `json:"text"`
	ParseMode         string `json:"parse_mode"`
}

// NewClient returns a new client to work with Telegram Bot API.
func NewClient(token string) *Client {
	return &Client{token: token}
}

// SendMessage sends a message to chat.
func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&sendMessage{
		ChatID:            chatID,
		DisableURLPreview: true,
		ParseMode:         "MarkdownV2",
		Text:              text,
	}); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", uri, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	result := &response{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println(string(result.Result))
	log.Println(string(result.Description))

	if !result.Ok {
		return &common.TGError{Code: result.ErrorCode, Description: string(result.Description)}
	}
	return nil
}
