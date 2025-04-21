package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WeChatNotifier struct {
	WebhookURL string
}

// WechatMessage 企业微信机器人 Markdown 消息结构体
type WechatMessage struct {
	MsgType  string         `json:"msgtype"`
	Markdown WechatMarkdown `json:"markdown"`
}
type WechatMarkdown struct {
	Content string `json:"content"`
}

func (w *WeChatNotifier) SendNotify(message, title string) error {
	msg := WechatMessage{
		MsgType: "markdown",
		Markdown: WechatMarkdown{
			Content: message,
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(w.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
