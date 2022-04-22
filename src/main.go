package main

import (
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Add handlers
	for _, handler := range handlers.Handlers {
		http.HandleFunc(handler.Endpoint, handler.Handler)
	}

	// Start server
	log.Info("Server listening")
	http.ListenAndServe(":3000", nil)
}
