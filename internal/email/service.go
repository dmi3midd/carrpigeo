package email

import (
	"bytes"
	"carrpigeo/internal/config"
	"carrpigeo/internal/htmltemplate"
	"context"
	"errors"
	"fmt"
	"html/template"
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
	// SendWithTemplate sends an email using a template.
	SendWithTemplate(ctx context.Context, to, subject, templateName string, data interface{}) error
}

type emailService struct {
	config       *config.SMTP
	client       EmailClient
	repository   EmailRepository
	templateRepo htmltemplate.HTMLTemplateRepository
}

func NewEmailService(client EmailClient, repository EmailRepository, templateRepo htmltemplate.HTMLTemplateRepository, cfg *config.SMTP) EmailService {
	return &emailService{
		config:       cfg,
		client:       client,
		repository:   repository,
		templateRepo: templateRepo,
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
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSendEmail, err)
	}

	if err := s.repository.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}

func (s *emailService) SendWithTemplate(ctx context.Context, to, subject, templateName string, data interface{}) error {
	op := "EmailService.SendWithTemplate"

	tmplData, err := s.templateRepo.GetByID(ctx, templateName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tmpl, err := template.New(templateName).Parse(tmplData.Content)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	email := Email{
		ID:       xid.New().String(),
		Sender:   s.config.User,
		Reciever: to,
		Subject:  subject,
		Body:     body.String(),
		IsHTML:   true,
		SentAt:   time.Now(),
	}

	if err := s.client.Send(&email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSendEmail, err)
	}

	if err := s.repository.Create(ctx, &email); err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToSaveEmail, err)
	}

	return nil
}
