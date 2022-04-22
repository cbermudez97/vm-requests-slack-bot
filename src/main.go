package main

import (
	"fmt"
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/handlers"
)

func main() {
	http.HandleFunc("/events-endpoint", handlers.Dispatcher)
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}
