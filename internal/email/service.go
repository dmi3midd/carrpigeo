package email

import (
	"carrpigeo/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	ErrTooManyRecipients = fmt.Errorf("too many recipients")
)

type EmailService interface {
	Send(ctx context.Context, to, subject, body string) error
	SendMany(ctx context.Context, to []string, subject, body string) error
}

type emailService struct {
	config     *config.SMTP
	client     EmailClient
	repository EmailRepository
}

func NewEmailService(client EmailClient, repository EmailRepository) EmailService {
	return &emailService{
		client:     client,
		repository: repository,
	}
}

func (s *emailService) Send(ctx context.Context, to, subject, body string) error {
	op := "EmailService.Send"
	email := Email{
		ID:       xid.New().String(),
		Sender:   s.config.User,
		Reciever: to,
		Subject:  subject,
		Body:     body,
		SentAt:   time.Now(),
	}
	if err := s.client.Send(&email); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repository.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *emailService) SendMany(ctx context.Context, to []string, subject, body string) error {
	op := "EmailService.SendMany"
	if len(to) > 10 {
		return fmt.Errorf("%s: %w", op, ErrTooManyRecipients)
	}
	emails := []Email{}
	for _, recipient := range to {
		email := Email{
			ID:       xid.New().String(),
			Sender:   s.config.User,
			Reciever: recipient,
			Subject:  subject,
			Body:     body,
			SentAt:   time.Now(),
		}
		emails = append(emails, email)
	}
	if err := s.client.SendMany(emails); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := s.repository.CreateMany(ctx, emails); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
