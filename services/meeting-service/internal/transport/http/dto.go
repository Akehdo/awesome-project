package http

import (
	"time"

	"meeting-service/internal/domain"
)

type createMeetingResponse struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	ObjectKey        string    `json:"object_key"`
	ContentType      string    `json:"content_type"`
	SizeBytes        int64     `json:"size_bytes"`
	CreatedAt        time.Time `json:"created_at"`
}

func newCreateMeetingResponse(
	meeting *domain.Meeting,
) createMeetingResponse {
	return createMeetingResponse{
		ID:               meeting.ID.String(),
		OriginalFilename: meeting.OriginalFilename,
		ObjectKey:        meeting.ObjectKey,
		ContentType:      meeting.ContentType,
		SizeBytes:        meeting.SizeBytes,
		CreatedAt:        meeting.CreatedAt,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
