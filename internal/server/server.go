package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"carrpigeo/internal/config"
	"carrpigeo/internal/postgres"
	"carrpigeo/internal/repository"
	"carrpigeo/internal/service"
)

type Server struct {
	db              postgres.PostgresService
	cfg             *config.Config
	emailService    service.EmailService
	templateService service.HTMLTemplateService
}

func NewServer(cfg *config.Config, db postgres.PostgresService) *http.Server {
	htmlTemplateRepository := repository.NewHTMLTemplateRepository(db.GetDB())
	htmlTemplateService := service.NewHTMLTemplateService(htmlTemplateRepository)

	emailRepository := repository.NewEmailRepository(db.GetDB())
	emailClient := service.NewEmailClient(&cfg.SMTP)
	emailService := service.NewEmailService(emailClient, emailRepository, htmlTemplateRepository, &cfg.SMTP)

	s := &Server{
		db:              db,
		cfg:             cfg,
		emailService:    emailService,
		templateService: htmlTemplateService,
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
