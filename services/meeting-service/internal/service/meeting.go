package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"meeting-service/internal/domain"
	"path"

	"github.com/google/uuid"
)

type ObjectStorage interface {
	Upload(
		ctx context.Context,
		objectKey string,
		reader io.Reader,
		size int64,
		contentType string,
	) error
	Delete(ctx context.Context, objectKey string) error
}

type MeetingRepository interface {
	Create(ctx context.Context, meeting *domain.Meeting) (*domain.Meeting, error)
}

var ErrMeetingFileRequired = errors.New("meeting file is required")

type CreateMeetingInput struct {
	OriginalFilename string
	ContentType      string
	SizeBytes        int64
	Reader           io.Reader
}

type MeetingService struct {
	repository MeetingRepository
	storage    ObjectStorage
}

func NewMeetingService(
	repository MeetingRepository,
	objectStorage ObjectStorage,
) *MeetingService {
	return &MeetingService{
		repository: repository,
		storage:    objectStorage,
	}
}

func (s *MeetingService) Create(
	ctx context.Context,
	input CreateMeetingInput,
) (*domain.Meeting, error) {
	if input.Reader == nil {
		return nil, ErrMeetingFileRequired
	}

	id := uuid.New()
	objectKey := path.Join("meetings", id.String(), "source")

	meeting, err := domain.NewMeeting(
		id,
		input.OriginalFilename,
		objectKey,
		input.ContentType,
		input.SizeBytes,
	)
	if err != nil {
		return nil, err
	}

	if err := s.storage.Upload(
		ctx,
		meeting.ObjectKey,
		input.Reader,
		meeting.SizeBytes,
		meeting.ContentType,
	); err != nil {
		return nil, fmt.Errorf("upload meeting file: %w", err)
	}

	createdMeeting, err := s.repository.Create(ctx, meeting)
	if err != nil {
		return nil, fmt.Errorf("create meeting: %w", err)
	}

	return createdMeeting, nil
}
