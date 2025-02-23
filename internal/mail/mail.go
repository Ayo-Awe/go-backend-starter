package mail

type Template string

type Mailer interface {
	Send(recipient string, msg EmailPayload) error
}

// An email message
type EmailPayload interface {
	Template() Template
}
