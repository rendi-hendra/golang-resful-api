package util

// Notifier defines the contract for sending notifications.
type Notifier interface {
	SendLoginNotification(to string) error
}
