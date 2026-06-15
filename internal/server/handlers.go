package server

import (
	errs "carrpigeo/internal/apierrors"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (s *Server) SendEmailHandler(w http.ResponseWriter, r *http.Request) error {
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

type CreateHTMLTemplateResponse struct {
	ID string `json:"id"`
}

func (s *Server) CreateHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(256 << 10); err != nil {
		return errs.NewInternalServerError(fmt.Errorf("Failed to parse multipart form: %w", err))
	}
	file, _, err := r.FormFile("file")
	name := r.FormValue("name")
	if err != nil {
		return errs.NewInternalServerError(fmt.Errorf("Failed to get file: %w", err))
	}
	if name == "" {
		return errs.NewBadRequestError(nil, "Name is required")
	}
	defer file.Close()

	ctx := r.Context()
	id, err := s.templateService.Save(ctx, name, &file)
	if err != nil {
		return errs.NewInternalServerError(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response := CreateHTMLTemplateResponse{ID: id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return errs.NewInternalServerError(err)
	}

	return nil
}

func (s *Server) RemoveHTMLTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	if id == "" {
		return errs.NewBadRequestError(nil, "ID is required")
	}

	ctx := r.Context()
	if err := s.templateService.Remove(ctx, id); err != nil {
		return errs.NewInternalServerError(err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

type SendEmailWithTemplateRequest struct {
	To         string      `json:"to"`
	Subject    string      `json:"subject"`
	TemplateID string      `json:"template_id"`
	Data       interface{} `json:"data"`
}

func (s *Server) SendEmailWithTemplateHandler(w http.ResponseWriter, r *http.Request) error {
	var req SendEmailWithTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.NewBadRequestError(err, "Failed to decode request")
	}
	defer r.Body.Close()

	if req.To == "" {
		return errs.NewBadRequestError(nil, "To is required")
	}
	if req.TemplateID == "" {
		return errs.NewBadRequestError(nil, "Template ID is required")
	}

	ctx := r.Context()
	if err := s.emailService.SendWithTemplate(ctx, req.To, req.Subject, req.TemplateID, req.Data); err != nil {
		return errs.NewInternalServerError(err)
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}
