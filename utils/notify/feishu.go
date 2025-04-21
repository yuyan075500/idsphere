package notify

import (
	"net/http"
	"strings"
	"time"
)

type FeishuNotifier struct {
	WebhookURL string
}

func (f *FeishuNotifier) SendNotify(message, title string) error {
	req, err := http.NewRequest("POST", f.WebhookURL, strings.NewReader(message))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
