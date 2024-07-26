package cdi

import (
	"fmt"
	"strings"

	"tags.cncf.io/container-device-interface/pkg/parser"
)

// Vendor represents the compatible CDI vendor.
type Vendor string

const (
	// Nvidia represents the Nvidia CDI vendor.
	Nvidia Vendor = "nvidia.com"
)

// ToVendor converts a string to a CDI vendor.
func ToVendor(vendor string) (Vendor, error) {
	switch vendor {
	case string(Nvidia):
		return Nvidia, nil
	default:
		return "", fmt.Errorf("invalid CDI vendor (%q)", vendor)
	}
}

// Class represents the compatible CDI class.
type Class string

const (
	// GPU is a single discrete GPU.
	GPU Class = "gpu"
	// IGPU is an integrated GPU.
	IGPU Class = "igpu"
	// MIG is a single MIG compatible GPU.
	MIG Class = "mig"
)

// ToClass converts a string to a CDI class.
func ToClass(c string) (Class, error) {
	switch c {
	case string(GPU):
		return GPU, nil
	case string(IGPU):
		return IGPU, nil
	case string(MIG):
		return MIG, nil
	default:
		return "", fmt.Errorf("invalid CDI class (%q)", c)
	}
}

// ID represents a Container Device Interface (CDI) identifier.
//
// +------------+-------+------------------------------+
// |   Vendor   | Class |           Name               |
// +---------------------------------------------------+
// | nvidia.com |  gpu  | [dev_idx] or `all`           |
// |            |  mig  | [dev_idx]:[mig_idx] or `all` |
// |            |  igpu | [dev_idx] or `all`           |
// +------------+-------+------------------------------+
//
// Examples:
//   - nvidia.com/gpu=0
//   - nvidia.com/gpu=all
//   - nvidia.com/mig=0:1
//   - nvidia.com/igpu=0
type ID struct {
	Vendor Vendor
	Class  Class
	Name   []string
}

// Empty returns true if the ID is empty.
func (id ID) Empty() bool {
	return id.Vendor == "" && id.Class == "" && len(id.Name) == 0
}

// ToCDI converts a string identifier to a CDI ID.
func ToCDI(id string) (ID, error) {
	vendor, class, name, err := parser.ParseQualifiedName(id)
	if err != nil {
		// The ID is not a valid CDI qualified name but it could be a valid DRM device ID.
		return ID{}, nil
	}

	vendorType, err := ToVendor(vendor)
	if err != nil {
		return ID{}, err
	}

	classType, err := ToClass(class)
	if err != nil {
		return ID{}, err
	}

	if name == "all" {
		return ID{Vendor: vendorType, Class: classType, Name: []string{"all"}}, nil
	}

	// The MIG nomenclature is specific to NVIDIA GPUs
	if classType == MIG && vendorType == Nvidia {
		// Try to split the name into <device>:<mig>
		migName := strings.Split(name, ":")
		if len(migName) != 2 {
			return ID{}, fmt.Errorf("invalid MIG CDI name (%q)", name)
		}

		return ID{Vendor: vendorType, Class: classType, Name: []string{migName[0], migName[1]}}, nil
	}

	return ID{Vendor: vendorType, Class: classType, Name: []string{name}}, nil
}
