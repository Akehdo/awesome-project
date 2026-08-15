package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	nethttp "net/http"

	"meeting-service/internal/domain"
	"meeting-service/internal/service"
)

const maxRequestBodySize = 1 << 20

type MeetingService interface {
	Create(
		ctx context.Context,
		input service.CreateMeetingInput,
	) (*domain.Meeting, error)
}

type MeetingHandler struct {
	service MeetingService
}

func NewMeetingHandler(service MeetingService) *MeetingHandler {
	return &MeetingHandler{
		service: service,
	}
}

func (h *MeetingHandler) Create(
	w nethttp.ResponseWriter,
	r *nethttp.Request,
) {
	var request createMeetingRequest

	r.Body = nethttp.MaxBytesReader(w, r.Body, maxRequestBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, nethttp.StatusBadRequest, "invalid request body")
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, nethttp.StatusBadRequest, "request body must contain one JSON object")
		return
	}

	meeting, err := h.service.Create(r.Context(), request.toServiceInput())
	if err != nil {
		if isCreateMeetingValidationError(err) {
			writeError(w, nethttp.StatusBadRequest, err.Error())
			return
		}

		log.Printf("create meeting: %v", err)
		writeError(w, nethttp.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, nethttp.StatusCreated, newCreateMeetingResponse(meeting))
}

func isCreateMeetingValidationError(err error) bool {
	return errors.Is(err, domain.ErrFilenameRequired) ||
		errors.Is(err, domain.ErrObjectKeyRequired) ||
		errors.Is(err, domain.ErrMeetingContentTypeEmpty) ||
		errors.Is(err, domain.ErrInvalidMeetingSize)
}
