package utils

import "github.com/selefra/selefra-provider-sdk/provider/schema"

func AddProviderPrefix(providerName string, d *schema.Diagnostics) *schema.Diagnostics {
	if IsEmpty(d) {
		return nil
	}
	newDiagnostics := schema.NewDiagnostics()
	for _, newD := range d.GetDiagnosticSlice() {

	}
	return newDiagnostics
}
