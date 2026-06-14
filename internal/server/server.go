package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"carrpigeo/internal/config"
	"carrpigeo/internal/database"
	"carrpigeo/internal/email"
	"carrpigeo/internal/htmltemplate"
)

type Server struct {
	db              database.DBService
	cfg             *config.Config
	emailService    email.EmailService
	templateService htmltemplate.HTMLTemplateService
}

func NewServer(cfg *config.Config, db database.DBService) *http.Server {
	htmlTemplateRepository := htmltemplate.NewHTMLTemplateRepository(db.GetDB())
	htmlTemplateService := htmltemplate.NewHTMLTemplateService(htmlTemplateRepository)

	emailRepository := email.NewEmailRepository(db.GetDB())
	emailClient := email.NewEmailClient(&cfg.SMTP)
	emailService := email.NewEmailService(emailClient, emailRepository, htmlTemplateRepository, &cfg.SMTP)

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
