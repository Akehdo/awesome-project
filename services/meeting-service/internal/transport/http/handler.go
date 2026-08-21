package http

import (
	"context"
	"errors"
	"io"
	"log"
	nethttp "net/http"

	"meeting-service/internal/domain"
	"meeting-service/internal/service"
)

const (
	maxFileSize       int64 = 100 << 20 // 100 МБ
	multipartOverhead       = 1 << 20   // запас на заголовки multipart
)

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
	r.Body = nethttp.MaxBytesReader(
		w,
		r.Body,
		maxFileSize+multipartOverhead,
	)
	file, header, err := r.FormFile("file")
	if err != nil {
		if _, ok := errors.AsType[*nethttp.MaxBytesError](err); ok {
			writeError(
				w,
				nethttp.StatusRequestEntityTooLarge,
				"file is too large",
			)
			return
		}

		writeError(w, nethttp.StatusBadRequest, "file is required")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	defer file.Close()

	if header.Size <= 0 {
		writeError(w, nethttp.StatusBadRequest, "file must not be empty")
		return
	}

	if header.Size > maxFileSize {
		writeError(
			w,
			nethttp.StatusRequestEntityTooLarge,
			"file is too large",
		)
		return
	}

	contentType, err := detectContentType(file)
	if err != nil {
		log.Printf("inspect uploaded file: %v", err)
		writeError(w, nethttp.StatusInternalServerError, "internal server error")
		return
	}

	if !isAllowedAudioContentType(contentType) {
		writeError(w, nethttp.StatusUnsupportedMediaType, "unsupported audio format")
		return
	}

	input := service.CreateMeetingInput{
		OriginalFilename: header.Filename,
		ContentType:      contentType,
		SizeBytes:        header.Size,
		Reader:           file,
	}

	meeting, err := h.service.Create(r.Context(), input)
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

func detectContentType(file io.ReadSeeker) (string, error) {
	buffer := make([]byte, 512)

	n, err := io.ReadFull(file, buffer)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return nethttp.DetectContentType(buffer[:n]), nil
}

func isAllowedAudioContentType(contentType string) bool {
	switch contentType {
	case "audio/mpeg", "audio/wave", "audio/aiff", "application/ogg":
		return true
	default:
		return false
	}
}

func isCreateMeetingValidationError(err error) bool {
	return errors.Is(err, service.ErrMeetingFileRequired) ||
		errors.Is(err, domain.ErrFilenameRequired) ||
		errors.Is(err, domain.ErrObjectKeyRequired) ||
		errors.Is(err, domain.ErrMeetingContentTypeEmpty) ||
		errors.Is(err, domain.ErrInvalidMeetingSize)
}
