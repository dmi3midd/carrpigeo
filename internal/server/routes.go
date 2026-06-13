package server

import (
	errs "carrpigeo/internal/errors"
	"encoding/json"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes

	mux.HandleFunc("GET /health", errs.ErrorHandler(s.healthHandler))

	mux.HandleFunc("POST /send/email", errs.ErrorHandler(s.sendEmailHandler))

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) error {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		return errs.NewInternalServerError(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}
