package htmltemplate

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
	ErrTemplateNotFound = errors.New("html template not found")
	ErrTemplateExists   = errors.New("html template already exists")
	ErrFailedToReadFile = errors.New("failed to read file")
)

type HTMLTemplateService interface {
	// Save saves html template in db.
	Save(ctx context.Context, name string, file *multipart.File) (string, error)
	// Remove removes html template from db.
	Remove(ctx context.Context, id string) error
	// GetByID returns html template by id.
	GetByID(ctx context.Context, id string) (*HTMLTemplate, error)
}

type htmlTemplateService struct {
	repo HTMLTemplateRepository
}

func NewHTMLTemplateService(repo HTMLTemplateRepository) HTMLTemplateService {
	return &htmlTemplateService{
		repo: repo,
	}
}

func (s *htmlTemplateService) Save(ctx context.Context, name string, file *multipart.File) (string, error) {
	op := "HTMLTemplateService.Create"

	contentBytes, err := io.ReadAll(*file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrFailedToReadFile)
	}

	id := xid.New().String()
	tmpl := &HTMLTemplate{
		ID:        id,
		Name:      name,
		Content:   string(contentBytes),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, tmpl); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *htmlTemplateService) Remove(ctx context.Context, id string) error {
	op := "HTMLTemplateService.Delete"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *htmlTemplateService) GetByID(ctx context.Context, id string) (*HTMLTemplate, error) {
	op := "HTMLTemplateService.GetByID"
	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tmpl, nil
}
