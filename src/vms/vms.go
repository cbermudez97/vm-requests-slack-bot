package vms

type VMRequest struct {
	Requester string
	Name      string
	Provider  string
	OS        string
	Type      string
	Region    string
	PrivateIP bool
}
