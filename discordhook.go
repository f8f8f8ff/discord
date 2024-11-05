package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		flagUrl      = flag.String("url", "", "required: discord hook url, also set from DISCORD_HOOK_URL")
		flagUsername = flag.String("username", "", "username override")
		flagContent  = flag.String("content", "", "required: message content, also read from stdin")
	)
	flag.Parse()
	url := *flagUrl
	if url == "" {
		url = os.Getenv("DISCORD_HOOK_URL")
		if url == "" {
			flag.Usage()
			os.Exit(1)
		}
	}
	content := *flagContent
	if content == "" {
		var err error
		content, err = readStdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		if content == "" {
			flag.Usage()
			os.Exit(1)
		}
	}
	m := message{
		url:      url,
		Content:  content,
		Username: *flagUsername,
	}

	err := m.Send()
	if err != nil {
		log.Fatal(err)
	}
}

type message struct {
	url      string
	Content  string `json:"content"`
	Username string `json:"username,omitempty"`
}

func (m message) Send() error {
	if m.url == "" {
		return fmt.Errorf("no url")
	}
	if m.Content == "" {
		return fmt.Errorf("no content")
	}
	messageJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, err := http.Post(m.url, "application/json", bytes.NewBuffer(messageJson))
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

func readStdinPipe() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("expected stdin from pipe")
	}
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
