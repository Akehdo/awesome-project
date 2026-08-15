package domain

import "errors"

var (
	ErrMeetingIDRequired       = errors.New("meeting id is required")
	ErrMeetingNotFound         = errors.New("meeting not found")
	ErrFilenameRequired        = errors.New("filename is required")
	ErrObjectKeyRequired       = errors.New("object key is required")
	ErrInvalidMeetingSize      = errors.New("meeting file size must be greater than zero")
	ErrMeetingContentTypeEmpty = errors.New("meeting content type is required")
)
