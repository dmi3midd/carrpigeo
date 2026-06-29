package service

import (
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/utils"
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
	ErrInvalidFileType  = errors.New("invalid file type")
)

type HTMLTemplateService interface {
	// Save saves html template in db.
	// Returns [ErrInvalidFileType] if file type is not text/html.
	Save(ctx context.Context, name string, file *multipart.File) (string, error)
	// Remove removes html template from db.
	Remove(ctx context.Context, id string) error
	// GetByID returns html template by id.
	// Returns [ErrTemplateNotFound] if template not found.
	GetByID(ctx context.Context, id string) (*domain.HTMLTemplate, error)
}

type htmlTemplateService struct {
	repo repository.HTMLTemplateRepository
}

func NewHTMLTemplateService(repo repository.HTMLTemplateRepository) HTMLTemplateService {
	return &htmlTemplateService{
		repo: repo,
	}
}

func (s *htmlTemplateService) Save(ctx context.Context, name string, file *multipart.File) (string, error) {
	op := "HTMLTemplateService.Create"

	mimeType, err := utils.DetectType(*file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if mimeType != "text/html" {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidFileType)
	}

	contentBytes, err := io.ReadAll(*file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id := xid.New().String()
	tmpl := &domain.HTMLTemplate{
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

func (s *htmlTemplateService) GetByID(ctx context.Context, id string) (*domain.HTMLTemplate, error) {
	op := "HTMLTemplateService.GetByID"
	tmpl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoTemplate) {
			return nil, fmt.Errorf("%s: %w", op, ErrTemplateNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tmpl, nil
}
