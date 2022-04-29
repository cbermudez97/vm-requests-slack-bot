package vms

import (
	"fmt"
	"strings"
)

const requestNotificationMsgTemplate = `VM Request

From: <@%s>
Name: %s
OS: %s
Provider: %s
Type: %s
Region: %s`

func haveAdditional(request VMRequest) bool {
	return request.PrivateIP
}

func BuildRequestNotificationMessage(request VMRequest) string {
	msg := fmt.Sprintf(
		requestNotificationMsgTemplate,
		request.Requester,
		request.Name,
		request.OS,
		request.Provider,
		request.Type,
		request.Region,
	)

	if haveAdditional(request) {
		withAdditional := []string{msg, "Additional:"}
		if request.PrivateIP {
			withAdditional = append(withAdditional, "Use Private Ip")
		}
		msg = strings.Join(withAdditional, "\n")
	}

	return msg
}
