package email

import (
	"carrpigeo/internal/config"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rs/xid"
)

var (
	ErrFailedToSendEmail = errors.New("failed to send email")
	ErrFailedToSaveEmail = errors.New("failed to save email")
)

type EmailService interface {
	// Send sends a single email.
	// Returns [ErrFailedToSendEmail] if failed to send email.
	// Returns [ErrFailedToSaveEmail] if failed to save email.
	Send(ctx context.Context, to, subject, body string) error
}

type emailService struct {
	config     *config.SMTP
	client     EmailClient
	repository EmailRepository
}

func NewEmailService(client EmailClient, repository EmailRepository, cfg *config.SMTP) EmailService {
	return &emailService{
		config:     cfg,
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

	go func() {
		if err := s.client.Send(&email); err != nil {
			slog.Error(
				"failed to send email",
				slog.String("id", email.ID),
				slog.String("to", email.Reciever),
				slog.Time("sent_at", email.SentAt),
				slog.String("error", err.Error()),
			)
		} else {
			slog.Info(
				"email sent successfully",
				slog.String("id", email.ID),
				slog.String("to", email.Reciever),
				slog.Time("sent_at", email.SentAt),
			)
		}
	}()

	if err := s.repository.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}
