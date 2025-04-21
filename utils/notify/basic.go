package notify

import "ops-api/model"

type Notifier interface {
	SendNotify(message, title string) error
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
