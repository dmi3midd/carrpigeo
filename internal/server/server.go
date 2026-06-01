package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"carrpigeo/internal/config"
	"carrpigeo/internal/database"
)

type Server struct {
	db  database.DBService
	cfg *config.Config
}

func NewServer(cfg *config.Config, db database.DBService) *http.Server {
	s := &Server{
		db:  db,
		cfg: cfg,
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
