package mail

import (
	"bytes"
	"text/template"

	"github.com/go-mail/mail/v2"
)

// concrete	email provider
type SMTPEmailProvider struct {
	dailer *mail.Dialer
	sender string
}

func NewSMTPEmailProvider(host string, port int, username, password, sender string) (*SMTPEmailProvider, error) {
	client := mail.NewDialer(host, port, username, password)

	return &SMTPEmailProvider{dailer: client, sender: sender}, nil
}

func (s *SMTPEmailProvider) Send(recipient string, data EmailPayload) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+string(data.Template()))
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", s.sender)
	msg.SetHeader("Subject", subject.String())

	// plain text should be added before the HTML part
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	return s.dailer.DialAndSend(msg)
}
