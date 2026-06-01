package email

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoEmail = errors.New("email not found")
)

type EmailRepository interface {
	GetByID(ctx context.Context, id string) (*Email, error)
	GetAll(ctx context.Context) ([]Email, error)
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

func (r *emailRepository) GetByID(ctx context.Context, id string) (*Email, error) {
	op := "EmailRepository.GetByID"
	query := `
	SELECT id, sender, reciever, subject, body, sent_at
	FROM emails
	WHERE id = $1
	`
	var email Email
	err := r.DB.GetContext(ctx, &email, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrNoEmail)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &email, nil
}

func (r *emailRepository) GetAll(ctx context.Context) ([]Email, error) {
	op := "EmailRepository.GetAll"
	query := `
	SELECT id, sender, reciever, subject, body, sent_at
	FROM emails
	`
	var emails []Email
	err := r.DB.SelectContext(ctx, &emails, query)
	if err != nil {
		return []Email{}, fmt.Errorf("%s: %w", op, err)
	}
	return emails, nil
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
