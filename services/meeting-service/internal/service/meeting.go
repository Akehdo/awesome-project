package service

import (
	"context"
	"fmt"
	"meeting-service/internal/domain"
)

type MeetingRepository interface {
	Create(ctx context.Context, meeting *domain.Meeting) (*domain.Meeting, error)
}
type CreateMeetingInput struct {
	OriginalFilename string
	ObjectKey        string
	ContentType      string
	SizeBytes        int64
}

type MeetingService struct {
	repository MeetingRepository
}

func NewMeetingService(repository MeetingRepository) *MeetingService {
	return &MeetingService{
		repository: repository,
	}
}

func (s *MeetingService) Create(
	ctx context.Context,
	input CreateMeetingInput,
) (*domain.Meeting, error) {
	meeting, err := domain.NewMeeting(
		input.OriginalFilename,
		input.ObjectKey,
		input.ContentType,
		input.SizeBytes,
	)
	if err != nil {
		return nil, err
	}

	createdMeeting, err := s.repository.Create(ctx, meeting)
	if err != nil {
		return nil, fmt.Errorf("create meeting: %w", err)
	}

	return createdMeeting, nil
}
