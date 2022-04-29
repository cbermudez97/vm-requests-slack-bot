package vms

import (
	"fmt"
	"regexp"

	"github.com/slack-go/slack"
)

var requesterRegex = regexp.MustCompile(`(?m:^From: <@(.*)>$)`)
var osRegex = regexp.MustCompile(`(?m:^OS: (.*)$)`)
var nameRegex = regexp.MustCompile(`(?m:^Name: (.*)$)`)
var typeRegex = regexp.MustCompile(`(?m:^Type: (.*)$)`)
var regionRegex = regexp.MustCompile(`(?m:^Region: (.*)$)`)
var providerRegex = regexp.MustCompile(`(?m:^Provider: (.*)$)`)
var privateIpRegex = regexp.MustCompile(`(?m:^Use Private Ip$)`)

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

func ParseRequestDataFrom(message slack.Message) (VMRequest, error) {
	data := VMRequest{}

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

		data.Requester, ok = apply(requesterRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.OS, ok = apply(osRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.Provider, ok = apply(providerRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.Name, ok = apply(nameRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.Region, ok = apply(regionRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.Type, ok = apply(typeRegex, requestRaw)
		if !ok {
			return data, fmt.Errorf("Invalid message format")
		}
		data.PrivateIP = privateIpRegex.Match([]byte(requestRaw))
	default:
		return data, fmt.Errorf("Invalid message structure")
	}

	return data, nil
}
