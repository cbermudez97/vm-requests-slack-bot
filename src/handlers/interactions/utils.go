package interactions

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

var requesterRegex = regexp.MustCompile(`(?m:^From: <@(.*)>$)`)
var distRegex = regexp.MustCompile(`(?m:^Distribution: (.*)$)`)
var vmTypeRegex = regexp.MustCompile(`(?m:^Type: (.*)$)`)

func apply(r *regexp.Regexp, s string) (string, bool) {
	matchs := r.FindAllStringSubmatch(s, 1)
	if len(matchs) < 1 {
		return "", false
	}
	groups := matchs[0]
	if len(groups) < 2 {
		return "", false
	}
	return groups[1], true
}

func parseRequestDataFrom(message slack.Message) (VMRequestData, error) {
	data := VMRequestData{}

	if len(message.Blocks.BlockSet) < 1 {
		return data, fmt.Errorf("Invalid message structure")
	}

	block := message.Blocks.BlockSet[0]

	switch block.BlockType() {
	case slack.MBTSection:
		sectionBlock, ok := block.(*slack.SectionBlock)
		if !ok {
			return data, fmt.Errorf("Invalid message structure")
		}
		requestRaw := sectionBlock.Text.Text
		log.Info(requestRaw)

		data.Requester, ok = apply(requesterRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.Distribution, ok = apply(distRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.Type, ok = apply(vmTypeRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}

	default:
		return data, fmt.Errorf("Invalid message structure")
	}

	return data, nil
}
