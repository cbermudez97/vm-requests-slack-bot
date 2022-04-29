package vms

type VMOS struct {
	Name  string
	Value string
}

var Ubuntu1804 = VMOS{
	Name:  "Ubuntu 18.04",
	Value: "ubuntu18.04",
}

var SupportedOS = [...]VMOS{
	Ubuntu1804,
}
