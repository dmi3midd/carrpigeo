package email

import (
	"carrpigeo/internal/config"
	"crypto/tls"
	"fmt"

	"gopkg.in/mail.v2"
)

type EmailClient interface {
	// Send sends a single email.
	Send(email *Email) error
	// SendMany sends multiple emails in a single SMTP session.
	SendMany(emails []Email) error
}

type emailClient struct {
	config *config.SMTP
	dialer *mail.Dialer
}

func NewSmtpClient(cfg *config.SMTP) EmailClient {
	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	dialer.TLSConfig = &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: true, // FALSE for production
	}

	return &emailClient{
		config: cfg,
		dialer: dialer,
	}
}

// buildMessage creates a new mail.Message from an Email struct.
func (c *emailClient) buildMessage(email *Email) *mail.Message {
	msg := mail.NewMessage()
	msg.SetHeader("From", c.config.User)
	msg.SetHeader("To", email.Reciever)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/plain", email.Body)
	return msg
}

func (c *emailClient) Send(email *Email) error {
	op := "SmtpClient.Send"

	msg := c.buildMessage(email)

	if err := c.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *emailClient) SendMany(emails []Email) error {
	op := "SmtpClient.SendMany"

	msgs := make([]*mail.Message, 0, len(emails))
	for _, e := range emails {
		msgs = append(msgs, c.buildMessage(&e))
	}

	if err := c.dialer.DialAndSend(msgs...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
