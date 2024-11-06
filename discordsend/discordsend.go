package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/f8f8f8ff/discord"
)

func main() {
	var (
		url       string = os.Getenv("DISCORD_HOOK_URL")
		content   string
		username  string = os.Getenv("DISCORD_USERNAME")
		status    string = os.Getenv("DISCORD_STATUS")
		title     string = os.Getenv("DISCORD_TITLE")
		ping      bool
		timestamp bool
	)

	flag.StringVar(&url, "url", url, "REQUIRED: discord webhook url [DISCORD_HOOK_URL]")
	flag.StringVar(&username, "username", username, "override username. [DISCORD_USERNAME]")
	flag.StringVar(&content, "content", content, "REQUIRED. also read from stdin")
	flag.StringVar(&status, "status", status, "one of info, warning, error. [DISCORD_STATUS]")
	flag.StringVar(&title, "title", title, "optional title for message. [DISCORD_TITLE]")
	flag.BoolVar(&ping, "ping", false, "ping @everyone")
	flag.BoolVar(&timestamp, "time", true, "include timestamp at start of message")
	flag.Parse()

	if url == "" {
		flag.Usage()
		os.Exit(1)
	}

	if content == "" {
		var err error
		content, err = readStdinPipe()
		if err != nil && err != ErrorStdinNoPipe {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			flag.Usage()
			os.Exit(1)
		}
		if content == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	content = "```\n" + content + "\n```"
	emb := discord.Embed{Content: content, Title: title}

	switch status {
	case "":
		break
	case "error":
		status = "# ERROR"
		emb.Color = discord.Red
		ping = true
	case "warning":
		status = "## warning"
		emb.Color = discord.Yellow
	case "info":
		emb.Color = discord.Blue
		fallthrough
	default:
		status = "### " + status
	}

	var body string
	if timestamp {
		body += time.Now().Format(time.DateTime)
	}
	if ping {
		body += " @everyone"
	}
	body += "\n" + status

	m := discord.Message{
		Url:      url,
		Username: username,
	}

	if emb.Color == 0 && emb.Title == "" {
		body += "\n" + content
		m.Content = body
	} else {
		m.Content = body
		m.Embeds = []discord.Embed{emb}
	}

	err := m.Send()
	if err != nil {
		log.Fatal(err)
	}
}

var ErrorStdinNoPipe error = errors.New("expected stdin from pipe")

func readStdinPipe() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", ErrorStdinNoPipe
	}
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
