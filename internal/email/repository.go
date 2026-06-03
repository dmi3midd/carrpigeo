package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrFailedToCreateEmail = errors.New("failed to create email")
)

type EmailRepository interface {
	// Create creates email in db.
	// Returns [ErrFailedToCreateEmail] if failed to create email.
	Create(ctx context.Context, email *Email) error
}

type emailRepository struct {
	DB *sqlx.DB
}

func NewEmailRepository(db *sqlx.DB) EmailRepository {
	return &emailRepository{
		DB: db,
	}
}

func (r *emailRepository) Create(ctx context.Context, email *Email) error {
	op := "EmailRepository.Create"
	query := `
	INSERT INTO emails (id, sender, reciever, subject, body, sent_at)
	VALUES (:id, :sender, :reciever, :subject, :body, :sent_at)
	`
	_, err := r.DB.NamedExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrFailedToCreateEmail, err)
	}
	return nil
}
