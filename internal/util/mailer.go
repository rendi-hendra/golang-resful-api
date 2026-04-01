package util

import (
	"fmt"
	"net/smtp"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Mailer struct {
	Config *viper.Viper
	Log    *logrus.Logger
}

func NewMailer(config *viper.Viper, logger *logrus.Logger) *Mailer {
	return &Mailer{
		Config: config,
		Log:    logger,
	}
}

func (m *Mailer) SendLoginNotification(to string) error {
	host := m.Config.GetString("mail.host")
	port := m.Config.GetString("mail.port")
	from := "no-reply@golang-resful.api"
	subject := "Subject: Login Notification\n"
	body := "Your account has just been logged into."
	msg := []byte(subject + "\n" + body)

	addr := fmt.Sprintf("%s:%s", host, port)
	simulation := m.Config.GetBool("mail.simulation")

	if simulation {
		m.Log.Infof("[SIMULATION] Email sent to: %s", to)
		m.Log.Infof("[SIMULATION] Subject: %s", subject)
		m.Log.Infof("[SIMULATION] Body: %s", body)
		return nil
	}

	// Mailpit doesn't require authentication by default on port 1025
	err := smtp.SendMail(addr, nil, from, []string{to}, msg)
	if err != nil {
		m.Log.Errorf("Failed to send email to %s: %+v", to, err)
		return err
	}

	m.Log.Infof("Email notification sent to %s", to)
	return nil
}
