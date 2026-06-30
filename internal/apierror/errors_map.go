package apierror

import (
	"carrpigeo/internal/service"
	"errors"
)

var ErrorMap = map[error]func(err error) error{
	service.ErrFailedToSendEmail: func(err error) error {
		return NewInternalServerError(err)
	},
	service.ErrFailedToSaveEmail: func(err error) error {
		return NewInternalServerError(err)
	},
	service.ErrInvalidFileType: func(err error) error {
		return NewBadRequestError(err, "Invalid file type")
	},
	service.ErrTemplateNotFound: func(err error) error {
		return NewNotFoundError(err, "Template not found")
	},
}

func MapError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr APIError
	if errors.As(err, &apiErr) {
		return err
	}

	for serviceErr, mapFn := range ErrorMap {
		if errors.Is(err, serviceErr) {
			return mapFn(err)
		}
	}

	return NewInternalServerError(err)
}
