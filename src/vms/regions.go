package vms

type VMRegion struct {
	Name  string
	Value string
}

var USCentral = VMRegion{
	Name:  "US Central",
	Value: "us-central",
}

var SupportedRegions = [...]VMRegion{
	USCentral,
}
