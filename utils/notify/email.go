package notify

import (
	"github.com/wonderivan/logger"
	"ops-api/utils/mail"
	"strings"
)

type EmailNotifier struct {
	To string
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
