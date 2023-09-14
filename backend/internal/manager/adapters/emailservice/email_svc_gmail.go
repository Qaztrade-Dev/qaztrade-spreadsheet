package emailservice

import (
	"context"
	"crypto/tls"
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"gopkg.in/gomail.v2"
)

//go:embed text_template.txt
var textTemplate []byte

type EmailServiceGmail struct {
	email    string
	password string
	template *template.Template
	dialer   *gomail.Dialer
}

var _ domain.EmailService = (*EmailServiceGmail)(nil)

func NewEmailServiceGmail(email, password string) *EmailServiceGmail {
	templ := template.Must(template.ParseFiles("./templates/forgot.html"))

	dialer := gomail.NewDialer("webmail.p-s.kz", 587, email, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailServiceGmail{
		email:    email,
		password: password,
		template: templ,
		dialer:   dialer,
	}
}

func (s *EmailServiceGmail) SendNotice(ctx context.Context, toEmail, mailName, fileName string, FileReader io.Reader) error {

	tempDir, err := os.MkdirTemp("", "mails")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tempDir)

	file, err := os.Create(filepath.Join(tempDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(file, FileReader); err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", s.email)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", mailName)
	m.SetBody("text/plain", string(textTemplate))
	m.Attach(file.Name())

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
