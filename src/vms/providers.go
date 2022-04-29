package vms

type VMProvider struct {
	Name  string
	Value string
}

var LinodeProvider = VMProvider{
	Name:  "Linode",
	Value: "linode",
}

var SupportedProviders = [...]VMProvider{
	LinodeProvider,
}
