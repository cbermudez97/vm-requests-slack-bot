package vms

type VMType struct {
	Name     string
	Provider VMProvider
	Value    string
}

var LinodeG6Standard1 = VMType{
	Name:     "G6 Standard 1",
	Provider: LinodeProvider,
	Value:    "g6-standard-1",
}

var AllTypes = [...]VMType{
	LinodeG6Standard1,
}

func SupportedTypesForProvider(provider VMProvider) []VMType {
	n := len(AllTypes) / len(SupportedProviders)
	providerVMTypes := make([]VMType, 0, n)

	for _, vmType := range AllTypes {
		if vmType.Provider.Value == provider.Value {
			providerVMTypes = append(providerVMTypes, vmType)
		}
	}
	return providerVMTypes
}
