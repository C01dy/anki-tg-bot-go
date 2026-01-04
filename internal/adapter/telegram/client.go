package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

func createUrl(c *Client, method string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
}

func (c *Client) FetchUpdates(offset int, limit int) ([]Update, error) {
	u := createUrl(c, "getUpdates")

	q := u.Query()
	q.Add("offset", fmt.Sprint(offset))
	q.Add("limit", fmt.Sprint(limit))
	u.RawQuery = q.Encode()

	res, err := c.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer res.Body.Close()

	var data UpdateResponse

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("couldn't decode updates: %w", err)
	}

	if !data.Ok {
		return nil, fmt.Errorf("telegram api error")
	}

	return data.Result, nil
}

type SendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *Client) SendMessage(chatID int64, text string) error {
	u := createUrl(c, "sendMessage")
	req := SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	dataToSend, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := c.client.Post(u.String(), "application/json", bytes.NewBuffer(dataToSend))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", res.StatusCode)
	}

	return nil
}
