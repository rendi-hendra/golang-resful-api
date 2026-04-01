package test

import (
	"net"
	"testing"
	"time"

	"github.com/rendi-hendra/resful-api/internal/config"
	"github.com/rendi-hendra/resful-api/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestSendLoginNotification(t *testing.T) {
	vConfig := config.NewViper()
	logger := config.NewLogger(vConfig)

	host := vConfig.GetString("mail.host")
	port := vConfig.GetString("mail.port")
	simulation := vConfig.GetBool("mail.simulation")
	address := host + ":" + port

	if !simulation {
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			t.Skip("Mailpit is not running on " + address + " and simulation mode is OFF, skipping email test")
		}
		conn.Close()
	}

	// Inisialisasi Mailer util
	mailer := util.NewMailer(vConfig, logger)

	// Uji coba pengiriman email login notification
	err := mailer.SendLoginNotification("test_user_mailer@example.com")

	// Pastikan tidak ada error balikan (Skenario Sukses)
	assert.Nil(t, err, "Diharapkan tidak ada error saat mengirim email lewat Mailpit")
}
