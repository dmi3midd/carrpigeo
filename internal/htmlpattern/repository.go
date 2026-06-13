package htmlpattern

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoPattern = errors.New("no html pattern in repository")
)

type HTMLPatternRepository interface {
	// GetByID returns html pattern by id.
	// Returns [ErrNoPattern] if html pattern not found.
	GetByID(ctx context.Context, id string) (*HTMLPattern, error)
	// Create creates email in db.
	Create(ctx context.Context, pattern *HTMLPattern) error
	// Delete deletes email from db.
	Delete(ctx context.Context, id string) error
}

type hTMLPatternRepository struct {
	db *sqlx.DB
}

func NewHTMLPatternRepository(db *sqlx.DB) HTMLPatternRepository {
	return &hTMLPatternRepository{
		db: db,
	}
}

func (r *hTMLPatternRepository) GetByID(ctx context.Context, id string) (*HTMLPattern, error) {
	op := "HTMLPatternRepository.GetByID"
	var pattern HTMLPattern
	query := `
	SELECT id, name, content, created_at
	FROM html_patterns
	WHERE id = $1
	`
	err := r.db.GetContext(ctx, &pattern, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoPattern)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &pattern, nil
}

func (r *hTMLPatternRepository) Create(ctx context.Context, pattern *HTMLPattern) error {
	op := "HTMLPatternRepository.Create"
	query := `
	INSERT INTO html_patterns (id, name, content, created_at)
	VALUES (:id, :name, :content, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, pattern)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *hTMLPatternRepository) Delete(ctx context.Context, id string) error {
	op := "HTMLPatternRepository.Delete"
	query := `
	DELETE FROM html_patterns
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
