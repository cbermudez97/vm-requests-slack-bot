package vms

import (
	"github.com/slack-go/slack"
)

// Modal data
const RequestModalCallbackId = "request-modal"

// Header section

func createHeader() *slack.SectionBlock {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			"*Please fill out the request form below:*",
			false,
			false,
		),
		nil,
		nil,
	)
}

// Name block data
const VMNameBlockId = "vm_name"
const VMNameActionId = "VM_NAME"

func createNameBlock() *slack.InputBlock {
	return slack.NewInputBlock(
		VMNameBlockId,
		slack.NewTextBlockObject(
			slack.PlainTextType,
			"VM Name",
			false,
			false,
		),
		slack.NewPlainTextInputBlockElement(
			nil,
			VMNameActionId,
		),
	)
}

// Distribution block data
const OSBlockId = "vm_os"
const OSActionId = "VM_OS"

func createOSOptions() []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(SupportedOS))
	for _, option := range SupportedOS {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, option.Name, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(option.Value, optionText, nil))
	}
	return optionBlockObjects
}

func createOSBlock() *slack.SectionBlock {
	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			"*Select a Distribution*",
			false,
			false,
		),
		nil,
		slack.NewAccessory(
			slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				nil,
				OSActionId,
				createOSOptions()...,
			),
		),
	)
	section.BlockID = OSBlockId
	return section
}

// VM provider block data
const VMProviderBlockId = "vm_provider"
const VMProviderActionId = "VM_PROVIDER"

func createProviderOptions() []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(SupportedProviders))
	for _, option := range SupportedProviders {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, option.Name, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(option.Value, optionText, nil))
	}
	return optionBlockObjects
}

func createProviderBlock() *slack.SectionBlock {
	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			"*Select a Provider*",
			false,
			false,
		),
		nil,
		slack.NewAccessory(
			slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				nil,
				VMProviderActionId,
				createProviderOptions()...,
			),
		),
	)
	section.BlockID = VMProviderBlockId
	return section
}

// VM type block data
const VMTypeBlockId = "vm_type"
const VMTypeActionId = "VM_TYPE"

func createTypeOptions(provider VMProvider) []*slack.OptionBlockObject {
	providerTypes := SupportedTypesForProvider(provider)
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(providerTypes))
	for _, option := range providerTypes {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, option.Name, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(option.Value, optionText, nil))
	}
	return optionBlockObjects
}

func createTypeBlock(provider VMProvider) *slack.SectionBlock {
	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			"*Select Type*",
			false,
			false,
		),
		nil,
		slack.NewAccessory(
			slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				nil,
				VMTypeActionId,
				createTypeOptions(provider)...,
			),
		),
	)
	section.BlockID = VMTypeBlockId
	return section
}

// VM region block data
const VMRegionBlockId = "vm_region"
const VMRegionActionId = "VM_region"

func createRegionOptions() []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(SupportedRegions))
	for _, option := range SupportedRegions {
		optionText := slack.NewTextBlockObject(slack.PlainTextType, option.Name, false, false)
		optionBlockObjects = append(optionBlockObjects, slack.NewOptionBlockObject(option.Value, optionText, nil))
	}
	return optionBlockObjects
}

func createRegionBlock() *slack.SectionBlock {
	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			"*Select a Region*",
			false,
			false,
		),
		nil,
		slack.NewAccessory(
			slack.NewOptionsSelectBlockElement(
				slack.OptTypeStatic,
				nil,
				VMRegionActionId,
				createRegionOptions()...,
			),
		),
	)
	section.BlockID = VMRegionBlockId
	return section
}

// VM Additional block data
const VMAdditionalBlockId = "vm_additional"
const VMAdditionalActionId = "VM_ADDITIONAL"

const VMPrivateIPValue = "USE_PRIVATE_IP"

func createPrivateIpBlock() *slack.OptionBlockObject {
	return slack.NewOptionBlockObject(
		VMPrivateIPValue,
		slack.NewTextBlockObject(
			slack.PlainTextType,
			"Use Private Ip",
			false,
			false,
		),
		nil,
	)
}

func createAdditionalInputBlock() *slack.SectionBlock {
	section := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.MarkdownType,
			"*Additional*",
			false,
			false,
		),
		nil,
		slack.NewAccessory(
			slack.NewCheckboxGroupsBlockElement(
				VMAdditionalActionId,
				createPrivateIpBlock(),
			),
		),
	)
	section.BlockID = VMAdditionalBlockId
	return section
}

func BuildVMRequestModal() slack.ModalViewRequest {
	// Modal texts
	titleText := slack.NewTextBlockObject(slack.PlainTextType, "VM Request", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Submit", false, false)
	// Header section
	headerSection := createHeader()
	// Name input
	vmNameBlock := createNameBlock()
	// OS input
	vmOSBlock := createOSBlock()
	// Provider input
	vmProviderBlock := createProviderBlock()
	// Region input
	vmRegionBlock := createRegionBlock()
	// Additional inputs
	vmAdditionalBlock := createAdditionalInputBlock()
	// Blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			slack.NewDividerBlock(),
			vmNameBlock,
			vmOSBlock,
			vmProviderBlock,
			vmRegionBlock,
			vmAdditionalBlock,
		},
	}
	// Modal
	var modalRequest slack.ModalViewRequest
	modalRequest.CallbackID = RequestModalCallbackId
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func BuildVMRequestModalWithTypes(provider VMProvider) slack.ModalViewRequest {
	// Modal texts
	titleText := slack.NewTextBlockObject(slack.PlainTextType, "VM Request", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Submit", false, false)
	// Header section
	headerSection := createHeader()
	// Name input
	vmNameBlock := createNameBlock()
	// OS input
	vmOSBlock := createOSBlock()
	// Provider input
	vmProviderBlock := createProviderBlock()
	// Type input
	vmTypeBlock := createTypeBlock(provider)
	// Region input
	vmRegionBlock := createRegionBlock()
	// Additional inputs
	vmAdditionalBlock := createAdditionalInputBlock()
	// Blocks
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			slack.NewDividerBlock(),
			vmNameBlock,
			vmOSBlock,
			vmProviderBlock,
			vmTypeBlock,
			vmRegionBlock,
			vmAdditionalBlock,
		},
	}
	// Modal
	var modalRequest slack.ModalViewRequest
	modalRequest.CallbackID = RequestModalCallbackId
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func NewVMRequestFromModal(i slack.InteractionCallback) VMRequest {
	stateValues := i.View.State.Values
	name := stateValues[VMNameBlockId][VMNameActionId].Value
	os := stateValues[OSBlockId][OSActionId].SelectedOption.Value
	provider := stateValues[VMProviderBlockId][VMProviderActionId].SelectedOption.Value
	vmType := ""
	vmTypeBlock, ok := stateValues[VMTypeBlockId]
	if ok {
		vmType = vmTypeBlock[VMTypeActionId].SelectedOption.Value
	}
	region := stateValues[VMRegionBlockId][VMRegionActionId].SelectedOption.Value

	privateIp := false
	for _, option := range stateValues[VMAdditionalBlockId][VMAdditionalActionId].SelectedOptions {
		if option.Value == VMPrivateIPValue {
			privateIp = true
		}
	}

	return VMRequest{
		Requester: i.User.ID,
		Name:      name,
		OS:        os,
		Provider:  provider,
		Type:      vmType,
		Region:    region,
		PrivateIP: privateIp,
	}
}
