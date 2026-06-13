package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"carrpigeo/internal/config"
	"carrpigeo/internal/database"
	"carrpigeo/internal/email"
	"carrpigeo/internal/htmlpattern"
)

type Server struct {
	db             database.DBService
	cfg            *config.Config
	emailService   email.EmailService
	patternService htmlpattern.HTMLPatternService
}

func NewServer(cfg *config.Config, db database.DBService) *http.Server {
	emailRepository := email.NewEmailRepository(db.GetDB())
	emailClient := email.NewEmailClient(&cfg.SMTP)
	emailService := email.NewEmailService(emailClient, emailRepository, &cfg.SMTP)

	htmlPatternRepository := htmlpattern.NewHTMLPatternRepository(db.GetDB())
	htmlPatternService := htmlpattern.NewHTMLPatternService(htmlPatternRepository)

	s := &Server{
		db:             db,
		cfg:            cfg,
		emailService:   emailService,
		patternService: htmlPatternService,
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
