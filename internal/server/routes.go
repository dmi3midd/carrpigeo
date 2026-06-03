package server

import (
	errs "carrpigeo/internal/errors"
	"encoding/json"
	"net/http"
	"slices"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes

	mux.HandleFunc("GET /health", errs.ErrorHandler(s.healthHandler))

	mux.HandleFunc("POST /send/email", errs.ErrorHandler(s.sendEmailHandler))

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := s.cfg.HTTPServer.CORS.AllowedOrigins

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
		}

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
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

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (s *Server) sendEmailHandler(w http.ResponseWriter, r *http.Request) error {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewBadRequestError(err, "Failed to decode request")
	}
	defer r.Body.Close()

	ctx := r.Context()
	if err := s.emailService.Send(ctx, req.To, req.Subject, req.Body); err != nil {
		return errs.NewInternalServerError(err)
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
