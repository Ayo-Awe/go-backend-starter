package testenv

import "github.com/ayo-awe/go-backend-starter/internal/mail"

type mockMailer struct {
	processedChan chan testMail
}

func newMockMailer() *mockMailer {
	return &mockMailer{
		processedChan: make(chan testMail, 3),
	}
}

type testMail struct {
	Recipient string
	Email     mail.EmailPayload
}

func (m *mockMailer) Send(recipient string, msg mail.EmailPayload) error {
	m.processedChan <- testMail{recipient, msg}
	return nil
}

// returns a buffered channel that returns processed mails
func (m *mockMailer) ProcessedChan() <-chan testMail {
	return m.processedChan
}

// removes all unprocessed mails from the mailer
func (m *mockMailer) Clear() {
	for range len(m.processedChan) {
		select {
		case <-m.processedChan:
		default:
		}
	}
}
