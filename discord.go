package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Url      string  `json:"-"`
	Content  string  `json:"content,omitempty"`
	Username string  `json:"username,omitempty"`
	Embeds   []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"description"`
	Color   int    `json:"color,omitempty"`
}

var (
	Red    int = 0xFF0000
	Green  int = 0x00FF00
	Blue   int = 0x0000FF
	Yellow int = 0xFFFF00
)

func (m Message) Send() error {
	// apply status

	err := m.Check()
	if err != nil {
		return err
	}
	messageJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, err := http.Post(m.Url, "application/json", bytes.NewBuffer(messageJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}
	return nil
}

func (m Message) Check() error {
	if m.Url == "" {
		return fmt.Errorf("no url")
	}
	if m.Content == "" && len(m.Embeds) == 0 {
		return fmt.Errorf("no content or embeds")
	}
	return nil
}
