package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Meeting struct {
	ID uuid.UUID

	// Исходное имя загруженного файла.
	// Например: daily-meeting.mp3
	OriginalFilename string

	// Путь к аудио внутри MinIO.
	// Например: meetings/{meeting_id}/daily-meeting.mp3
	ObjectKey string

	// MIME-тип файла.
	// Например: audio/mpeg
	ContentType string

	// Размер аудиофайла в байтах.
	SizeBytes int64

	CreatedAt time.Time
}

func NewMeeting(
	originalFilename string,
	objectKey string,
	contentType string,
	sizeBytes int64,
) (*Meeting, error) {
	originalFilename = strings.TrimSpace(originalFilename)
	objectKey = strings.TrimSpace(objectKey)
	contentType = strings.TrimSpace(contentType)

	if originalFilename == "" {
		return nil, ErrFilenameRequired
	}

	if objectKey == "" {
		return nil, ErrObjectKeyRequired
	}

	if contentType == "" {
		return nil, ErrMeetingContentTypeEmpty
	}

	if sizeBytes <= 0 {
		return nil, ErrInvalidMeetingSize
	}

	return &Meeting{
		OriginalFilename: originalFilename,
		ObjectKey:        objectKey,
		ContentType:      contentType,
		SizeBytes:        sizeBytes,
		CreatedAt:        time.Now().UTC(),
	}, nil
}
