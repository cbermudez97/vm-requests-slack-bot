package handlers

import "net/http"

type Handler struct {
	Endpoint string
	Handler  func(w http.ResponseWriter, r *http.Request)
}

var Handlers = [...]Handler{
	DispatcherHandler,
	RequestHandler,
	InteractionHandler,
}
