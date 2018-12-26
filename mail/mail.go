package mail

import (
	"crypto/tls"

	//"github.com/natefinch/lumberjack"
	gomail "gopkg.in/gomail.v2"
)

type Mailer struct {
	Addr       string
	Port       string
	User       string
	Password   string
	Skipverify bool
}

func New(addr, port, user, pass string) *Mailer {
	return &Mailer{
		Addr:       addr,
		Port:       port,
		User:       user,
		Password:   pass,
		Skipverify: true,
	}
}

func (m *Mailer) Mail(from, subject, body string, receivers []string) error {
	mm := gomail.NewMessage()
	mm.SetHeader("From", from)
	mm.SetHeader("To", receivers...)
	mm.SetHeader("Subject", subject)
	mm.SetBody("text/plain", body)
	d := gomail.NewDialer(m.Addr, m.Port, m.User, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: m.Skipverify}

	return d.DialAndSend(mm)
}
