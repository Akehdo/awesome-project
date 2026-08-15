package http

import nethttp "net/http"

func NewRouter(meetingHandler *MeetingHandler) nethttp.Handler {
	router := nethttp.NewServeMux()
	router.HandleFunc("POST /meetings", meetingHandler.Create)

	return router
}
