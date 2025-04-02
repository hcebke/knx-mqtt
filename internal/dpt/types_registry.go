package dpt

import (
	"reflect"

	"github.com/vapourismo/knx-go/knx/dpt"
)

var dptTypes = map[string]interface{}{
	// 19.xxx
	"19.001": new(DPT_19001),
}

// Produce returns a new instance of the specified datapoint type.
// It first checks if the datapoint type is registered in the local registry,
// and if not, it falls back to the knx-go library's registry.
func Produce(name string) (dpt.Datapoint, bool) {
	// First, check if the datapoint type is registered in the local registry
	if x, ok := dptTypes[name]; ok {
		d_type := reflect.TypeOf(x).Elem()
		d := reflect.New(d_type).Interface()

		// Convert to dpt.Datapoint interface
		if dp, ok := d.(dpt.Datapoint); ok {
			return dp, true
		}
	}

	// If not found in the local registry, fall back to the knx-go library's registry
	return dpt.Produce(name)
}

// ListSupportedTypes returns the names of all known datapoint types (DPTs).
// It combines the types from the local registry and the knx-go library's registry.
func ListSupportedTypes() []string {
	// Get the types from the knx-go library
	knxGoTypes := dpt.ListSupportedTypes()

	// Add the types from the local registry
	for k := range dptTypes {
		// Check if the type is already in the list
		found := false
		for _, t := range knxGoTypes {
			if t == k {
				found = true
				break
			}
		}

		// If not found, add it to the list
		if !found {
			knxGoTypes = append(knxGoTypes, k)
		}
	}

	return knxGoTypes
}
