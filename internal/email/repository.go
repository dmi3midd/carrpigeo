package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoEmail = errors.New("email not found")
)

type EmailRepository interface {
	Create(ctx context.Context, email *Email) error
	CreateMany(ctx context.Context, emails []Email) error
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
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *emailRepository) CreateMany(ctx context.Context, emails []Email) error {
	op := "EmailRepository.CreateMany"
	query := `
	INSERT INTO emails (id, sender, reciever, subject, body, sent_at)
	VALUES (:id, :sender, :reciever, :subject, :body, :sent_at)
	`
	_, err := r.DB.NamedExecContext(ctx, query, emails)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
