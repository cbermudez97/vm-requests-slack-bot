package interactions

import "github.com/slack-go/slack"

func parseRequestDataFrom(message slack.Message) (VMRequestData, error) {
	data := VMRequestData{}

	// FIXME: fix parsing
	// if len(message.Blocks.BlockSet) < 1 {
	// 	return data, fmt.Errorf("Invalid message structure")
	// }

	// msgTextSection := message.Blocks.BlockSet[0]
	// switch msgTextSection.(type) {
	// case slack.SectionBlock:
	// 	// TODO: parse data from text
	// default:
	// 	return data, fmt.Errorf("Invalid message structure")
	// }

	return data, nil
}
