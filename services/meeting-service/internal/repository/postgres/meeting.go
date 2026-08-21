package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"meeting-service/internal/domain"
)

const (
	createMeetingQuery = `
		INSERT INTO meetings (
		    id,
			original_filename,
			object_key,
			content_type,
			size_bytes,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
)

type MeetingRepository struct {
	db *sql.DB
}

func NewMeetingRepository(db *sql.DB) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) Create(ctx context.Context, meeting *domain.Meeting) (*domain.Meeting, error) {
	if meeting == nil {
		return nil, errors.New("meeting is required")
	}

	err := r.db.QueryRowContext(
		ctx,
		createMeetingQuery,
		meeting.ID,
		meeting.OriginalFilename,
		meeting.ObjectKey,
		meeting.ContentType,
		meeting.SizeBytes,
		meeting.CreatedAt,
	).Scan(&meeting.ID)
	if err != nil {
		return nil, fmt.Errorf("create meeting: %w", err)
	}

	return meeting, nil
}
