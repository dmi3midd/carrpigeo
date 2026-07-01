package server

import (
	"net/http"

	"carrpigeo/internal/config"
	"carrpigeo/internal/postgres"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/service"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	postgres        postgres.PostgresService
	cfg             *config.Config
	emailService    service.EmailService
	templateService service.HTMLTemplateService
	validator       *validator.Validate
}

func NewServer(cfg *config.Config, postgres postgres.PostgresService) *http.Server {
	validator := validator.New()
	htmlTemplateRepository := repository.NewHTMLTemplateRepository(postgres.GetDB())
	htmlTemplateService := service.NewHTMLTemplateService(htmlTemplateRepository)

	emailRepository := repository.NewEmailRepository(postgres.GetDB())
	emailClient := service.NewEmailClient(&cfg.SMTP)
	emailService := service.NewEmailService(emailClient, emailRepository, htmlTemplateRepository, &cfg.SMTP)

	s := &Server{
		postgres:        postgres,
		cfg:             cfg,
		emailService:    emailService,
		templateService: htmlTemplateService,
		validator:       validator,
	}

	router := s.RegisterRoutes()
	return &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}
}
