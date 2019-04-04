//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pcidb

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	cacheOnlyTrue = true
	localOnlyTrue = true
)

// ProgrammingInterface is the PCI programming interface for a class of PCI
// devices
type ProgrammingInterface struct {
	// hex-encoded PCI_ID of the programming interface
	ID string `json:"id"`
	// common string name for the programming interface
	Name string `json:"name"`
}

// Subclass is a subdivision of a PCI class
type Subclass struct {
	// hex-encoded PCI_ID for the device subclass
	ID string `json:"id"`
	// common string name for the subclass
	Name string `json:"name"`
	// any programming interfaces this subclass might have
	ProgrammingInterfaces []*ProgrammingInterface `json:"programming_interfaces"`
}

// Class is the PCI class
type Class struct {
	// hex-encoded PCI_ID for the device class
	ID string `json:"id"`
	// common string name for the class
	Name string `json:"name"`
	// any subclasses belonging to this class
	Subclasses []*Subclass `json:"subclasses"`
}

// Product provides information about a PCI device model
// NOTE(jaypipes): In the hardware world, the PCI "device_id" is the identifier
// for the product/model
type Product struct {
	// vendor ID for the product
	VendorID string `json:"vendor_id"`
	// hex-encoded PCI_ID for the product/model
	ID string `json:"id"`
	// common string name of the vendor
	Name string `json:"name"`
	// "subdevices" or "subsystems" for the product
	Subsystems []*Product `json:"subsystems"`
}

// Vendor provides information about a device vendor
type Vendor struct {
	// hex-encoded PCI_ID for the vendor
	ID string `json:"id"`
	// common string name of the vendor
	Name string `json:"name"`
	// all top-level devices for the vendor
	Products []*Product `json:"products"`
}

type PCIDB struct {
	// hash of class ID -> class information
	Classes map[string]*Class `json:"classes"`
	// hash of vendor ID -> vendor information
	Vendors map[string]*Vendor `json:"vendors"`
	// hash of vendor ID + product/device ID -> product information
	Products map[string]*Product `json:"products"`
}

// WithOption is used to represent optionally-configured settings
type WithOption struct {
	// Chroot is the directory that pcidb uses when attempting to discover
	// pciids DB files
	Chroot *string
	// CacheOnly is mostly just useful for testing. It essentially disables
	// looking for any non ~/.cache/pci.ids filepaths (which is useful when we
	// want to test the fetch-from-network code paths
	CacheOnly *bool
	// Path provides a search path directly to find pci.ids or pci.ids
	Path *string
	// LocalOnly disables any fetch-from-network capability, for use in
	// environments were it is undesirable
	LocalOnly *bool
}

func WithChroot(dir string) *WithOption {
	return &WithOption{Chroot: &dir}
}

func WithCacheOnly() *WithOption {
	return &WithOption{CacheOnly: &cacheOnlyTrue}
}

func WithPath(path string) *WithOption {
	return &WithOption{Path: &path}
}

func WithLocalOnly() *WithOption {
	return &WithOption{LocalOnly: &localOnlyTrue}
}

func mergeOptions(opts ...*WithOption) *WithOption {
	// Grab options from the environs by default
	defaultChroot := "/"
	if val, exists := os.LookupEnv("PCIDB_CHROOT"); exists {
		defaultChroot = val
	}
	defaultCacheOnly := false
	if val, exists := os.LookupEnv("PCIDB_CACHE_ONLY"); exists {
		if parsed, err := strconv.ParseBool(val); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Failed parsing a bool from PCIDB_CACHE_ONLY "+
					"environ value of %s",
				val,
			)
		} else if parsed {
			defaultCacheOnly = parsed
		}
	}
	defaultPath := filepath.Join(defaultChroot, "usr", "share", "hwdata", "pci.ids")
	if val, exists := os.LookupEnv("PCIDB_PATH"); exists {
		defaultChroot = "/"
		defaultPath = val
	}
	defaultLocalOnly := false
	if val, exists := os.LookupEnv("PCIDB_LOCAL_ONLY"); exists {
		if parsed, err := strconv.ParseBool(val); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Failed parsing a bool from PCIDB_LOCAL_ONLY "+
					"environ value of %s",
				val,
			)
		} else if parsed {
			defaultLocalOnly = parsed
		}
	}

	merged := &WithOption{}
	for _, opt := range opts {
		if opt.Chroot != nil {
			merged.Chroot = opt.Chroot
		}
		if opt.CacheOnly != nil {
			merged.CacheOnly = opt.CacheOnly
		}
		if opt.Path != nil {
			merged.Path = opt.Path
		}
		if opt.LocalOnly != nil {
			merged.LocalOnly = opt.LocalOnly
		}
	}
	// Set the default value if missing from merged
	if merged.Chroot == nil {
		merged.Chroot = &defaultChroot
	}
	if merged.CacheOnly == nil {
		merged.CacheOnly = &defaultCacheOnly
	}
	if merged.Path == nil {
		merged.Path = &defaultPath
	}
	if merged.LocalOnly == nil {
		merged.LocalOnly = &defaultLocalOnly
	}

	return merged
}

// New returns a pointer to a PCIDB struct which contains information you can
// use to query PCI vendor, product and class information. It accepts zero or
// more pointers to WithOption structs. If you want to modify the behaviour of
// pcidb, use one of the option modifiers when calling New. For example, to
// change the root directory that pcidb uses when discovering pciids DB files,
// call New(WithChroot("/my/root/override"))
func New(opts ...*WithOption) (*PCIDB, error) {
	ctx := contextFromOptions(mergeOptions(opts...))
	db := &PCIDB{}
	err := db.load(ctx)
	return db, err
}
