package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

var from = os.Getenv("EMAIL_SENDER")
var user = os.Getenv("EMAIL_USER")
var password = os.Getenv("EMAIL_PASSWORD")

type Mailer struct {
	To      []string
	Subject string
	// relative to the current directory
	TemplatePath string
	Variables    map[string]string
}

func (m *Mailer) SendWithTemplate() error {
	// TODO proper handling depending on the environment
	// auth := smtp.PlainAuth("", user, password, smtpHost)

	t, err := template.ParseFiles(m.TemplatePath)
	if err != nil {
		//TODO proper error handling
		return err
	}

	body := new(bytes.Buffer)

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n", m.Subject, mimeHeaders)))

	t.Execute(body, m.Variables)

	err = smtp.SendMail(os.Getenv("EMAIL_ADDR"), nil, os.Getenv("EMAIL_SENDER"), m.To, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
