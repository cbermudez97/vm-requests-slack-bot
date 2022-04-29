package handlers

import (
	"net/http"

	"github.com/cbermudez97/vm-requests-slack-bot/src/vms"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Handler configs
const requestsEndpoint = "/request-vm"
const requestCommandStr = "/request-vm"

func requestsCreation(w http.ResponseWriter, r *http.Request) {
	// Verify signing secret
	if err := VerifySigningSecret(r); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Build command
	rCmd, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Handle command
	if rCmd.Command == requestCommandStr {
		api := GetApi()

		modalRequest := vms.BuildVMRequestModal()
		_, err := api.OpenView(rCmd.TriggerID, modalRequest)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Errorf("Invalid command executed. Expected %s but got %s", requestCommandStr, rCmd.Command)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

var RequestHandler = Handler{
	Endpoint: requestsEndpoint,
	Handler:  requestsCreation,
}
