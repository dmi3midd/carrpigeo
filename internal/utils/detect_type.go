package utils

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

func DetectType(file multipart.File) (string, error) {
	op := "utils_package.DetectType"
	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)
	mimeType := http.DetectContentType(buffer)
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("%s: failed to seek to start: %w", op, err)
	}
	return mimeType, nil
}
