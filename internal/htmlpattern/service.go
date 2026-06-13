package htmlpattern

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/rs/xid"
)

var (
	ErrPatternNotFound  = errors.New("html pattern not found")
	ErrPatternExists    = errors.New("html pattern already exists")
	ErrFailedToReadFile = errors.New("failed to read file")
)

type HTMLPatternService interface {
	// Save saves html pattern in db.
	Save(ctx context.Context, name string, file *multipart.File) (string, error)
	// Remove removes html pattern from db.
	Remove(ctx context.Context, id string) error
}

type htmlPatternService struct {
	repo HTMLPatternRepository
}

func NewHTMLPatternService(repo HTMLPatternRepository) HTMLPatternService {
	return &htmlPatternService{
		repo: repo,
	}
}

func (s *htmlPatternService) Save(ctx context.Context, name string, file *multipart.File) (string, error) {
	op := "HTMLPatternService.Create"

	contentBytes, err := io.ReadAll(*file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrFailedToReadFile)
	}

	id := xid.New().String()
	pattern := &HTMLPattern{
		ID:        id,
		Name:      name,
		Content:   string(contentBytes),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, pattern); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *htmlPatternService) Remove(ctx context.Context, id string) error {
	op := "HTMLPatternService.Delete"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
