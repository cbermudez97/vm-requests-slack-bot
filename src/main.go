package main

import (
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers/interactions"
	log "github.com/sirupsen/logrus"
)

var Handlers = [...]handlers.Handler{
	handlers.DispatcherHandler,
	handlers.RequestHandler,
	interactions.InteractionHandler,
}

func main() {
	// Add handlers
	for _, handler := range Handlers {
		http.HandleFunc(handler.Endpoint, handler.Handler)
	}

	// Start server
	log.Info("Server listening")
	http.ListenAndServe(":3000", nil)
}
