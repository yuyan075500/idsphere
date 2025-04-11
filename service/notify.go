package service

import (
	"bytes"
	"encoding/json"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/model"
	"ops-api/utils/mail"
	"strings"
)

type Notifier interface {
	SendNotify(message, title string) error
}

type EmailNotifier struct {
	To string
}
type WeChatNotifier struct {
	WebhookURL string
}
type DingTalkNotifier struct {
	WebhookURL string
}
type FeishuNotifier struct {
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

func (e *EmailNotifier) SendNotify(message, title string) error {
	receivers := strings.Split(e.To, ",")
	for _, receiver := range receivers {
		if err := mail.Email.SendMsg([]string{receiver}, nil, nil, title, message, "html"); err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return nil
}

func (w *WeChatNotifier) SendNotify(message, title string) error {
	return nil
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

func (f *FeishuNotifier) SendNotify(message, title string) error {
	return nil
}

// GetNotifier 获取通知类型
func GetNotifier(task model.ScheduledTask) Notifier {

	var (
		notifyType = *task.NotifyType
		receiver   = *task.Receiver
	)

	switch notifyType {
	case 1:
		return &EmailNotifier{To: receiver}
	case 2:
		return &DingTalkNotifier{WebhookURL: receiver}
	case 3:
		return &FeishuNotifier{WebhookURL: receiver}
	case 4:
		return &WeChatNotifier{WebhookURL: receiver}
	default:
		return nil
	}
}
