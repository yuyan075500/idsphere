package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type DingTalkNotifier struct {
	WebhookURL string
}

// DingTalkMessage 钉钉机器人 Markdown 消息结构体
type DingTalkMessage struct {
	MsgType  string           `json:"msgtype"`
	Markdown DingTalkMarkdown `json:"markdown"`
}
type DingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (d *DingTalkNotifier) SendNotify(message, title string) error {
	msg := DingTalkMessage{
		MsgType: "markdown",
		Markdown: DingTalkMarkdown{
			Title: title,
			Text:  message,
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
