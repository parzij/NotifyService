package models

type NotificationChannel string
type NotificationStatus string

const (
	ChannelEmail    NotificationChannel = "email"
	ChannelTelegram NotificationChannel = "telegram"
	ChannelSMS      NotificationChannel = "sms"
)

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

func (c NotificationChannel) isValid() bool {
	switch c {
	case ChannelEmail, ChannelTelegram, ChannelSMS:
		return true
	}
	return false
}

func (c NotificationStatus) isValid() bool {
	switch c {
	case StatusPending, StatusSent, StatusFailed:
		return true
	}
	return false
}
