package util

// INotifier defines the contract for sending notifications.
type INotifier interface {
	SendLoginNotification(to string) error
}
