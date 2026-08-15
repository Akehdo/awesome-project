package http

import (
	"time"

	"meeting-service/internal/domain"
	"meeting-service/internal/service"
)

type createMeetingRequest struct {
	OriginalFilename string `json:"original_filename"`
	ObjectKey        string `json:"object_key"`
	ContentType      string `json:"content_type"`
	SizeBytes        int64  `json:"size_bytes"`
}

type createMeetingResponse struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	ObjectKey        string    `json:"object_key"`
	ContentType      string    `json:"content_type"`
	SizeBytes        int64     `json:"size_bytes"`
	CreatedAt        time.Time `json:"created_at"`
}

func (r createMeetingRequest) toServiceInput() service.CreateMeetingInput {
	return service.CreateMeetingInput{
		OriginalFilename: r.OriginalFilename,
		ObjectKey:        r.ObjectKey,
		ContentType:      r.ContentType,
		SizeBytes:        r.SizeBytes,
	}
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
