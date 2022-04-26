package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Handler configs
const requestsEndpoint = "/request-vm"
const requestCommandStr = "/request-vm"

// Modal data
const requestModalCallbackId = "request-modal"

// Distribution block data
const DistBlockId = "dist"
const DistActionId = "DIST"

var DistOptions = []string{
	"Ubuntu 20.04",
	"Ubunut 18.04",
}

// VM type block data
const VMTypeBlockId = "vm_type"
const VMTypeActionId = "VM_TYPE"

var VMTypeOptions = []string{
	"Linode 1",
	"Linode 2",
}

func createOptionBlockObjects(options []string) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(options))
	for _, option := range options {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, option, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(option, optionText, nil))
	}
	return optionBlockObjects
}

func buildVMRequestModal() slack.ModalViewRequest {
	// Modal texts
	titleText := slack.NewTextBlockObject(slack.PlainTextType, "VM Request", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Submit", false, false)

	// Header section
	headerText := slack.NewTextBlockObject(slack.MarkdownType, "*Please fill out the request form below:*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Distribution input
	distOptionsElems := createOptionBlockObjects(DistOptions)
	distText := slack.NewTextBlockObject(slack.PlainTextType, "Select a Distribution", false, false)
	distOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, DistActionId, distOptionsElems...)
	distBlock := slack.NewInputBlock(DistBlockId, distText, distOption)

	// VM Type input
	vmTypeOptionsElems := createOptionBlockObjects(VMTypeOptions)
	vmTypeText := slack.NewTextBlockObject(slack.PlainTextType, "Select a VM Type", false, false)
	vmTypeOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic, nil, VMTypeActionId, vmTypeOptionsElems...)
	vmTypeBlock := slack.NewInputBlock(VMTypeBlockId, vmTypeText, vmTypeOption)

	// Additional details
	// TODO: define additional details modal

	// Blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			distBlock,
			vmTypeBlock,
		},
	}

	// Modal
	var modalRequest slack.ModalViewRequest
	modalRequest.CallbackID = requestModalCallbackId
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

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

		modalRequest := buildVMRequestModal()
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
