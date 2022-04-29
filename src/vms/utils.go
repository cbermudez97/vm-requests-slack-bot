package vms

import (
	"fmt"
)

func FindProviderByValue(providerValue string) (*VMProvider, error) {
	for _, provider := range SupportedProviders {
		if provider.Value == providerValue {
			return &provider, nil
		}
	}

	return nil, fmt.Errorf("No provider found")
}
