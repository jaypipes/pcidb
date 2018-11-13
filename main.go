//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pcidb

import (
	"fmt"
	"os"
	"strconv"
)

var (
	cacheOnlyTrue = true
)

type PCIProgrammingInterface struct {
	// Id is DEPRECATED in 0.2 and will be removed in the 1.0 release. Please
	// use the equivalent ID field.
	Id   string
	ID   string // hex-encoded PCI_ID of the programming interface
	Name string // common string name for the programming interface
}

type PCISubclass struct {
	// Id is DEPRECATED in 0.2 and will be removed in the 1.0 release. Please
	// use the equivalent ID field.
	Id                    string
	ID                    string                     // hex-encoded PCI_ID for the device subclass
	Name                  string                     // common string name for the subclass
	ProgrammingInterfaces []*PCIProgrammingInterface // any programming interfaces this subclass might have
}

type PCIClass struct {
	// Id is DEPRECATED in 0.2 and will be removed in the 1.0 release. Please
	// use the equivalent ID field.
	Id         string
	ID         string         // hex-encoded PCI_ID for the device class
	Name       string         // common string name for the class
	Subclasses []*PCISubclass // any subclasses belonging to this class
}

// NOTE(jaypipes): In the hardware world, the PCI "device_id" is the identifier
// for the product/model
type PCIProduct struct {
	// VendorId is DEPRECATED in 0.2 and will be removed in the 1.0 release. Please
	// use the equivalent VendorID field.
	VendorId string
	VendorID string // vendor ID for the product
	// Id is DEPRECATED in 0.2 and will be removed in the 1.0 release. Please
	// use the equivalent ID field.
	Id         string
	ID         string        // hex-encoded PCI_ID for the product/model
	Name       string        // common string name of the vendor
	Subsystems []*PCIProduct // "subdevices" or "subsystems" for the product
}

type PCIVendor struct {
	// Id is DEPRECATED in 0.2 and will be removed in the 1.0 release. Please
	// use the equivalent ID field.
	Id       string
	ID       string        // hex-encoded PCI_ID for the vendor
	Name     string        // common string name of the vendor
	Products []*PCIProduct // all top-level devices for the vendor
}

type PCIDB struct {
	// hash of class ID -> class dbrmation
	Classes map[string]*PCIClass
	// hash of vendor ID -> vendor dbrmation
	Vendors map[string]*PCIVendor
	// hash of vendor ID + product/device ID -> product dbrmation
	Products map[string]*PCIProduct
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
}

func WithChroot(dir string) *WithOption {
	return &WithOption{Chroot: &dir}
}

func WithCacheOnly() *WithOption {
	return &WithOption{CacheOnly: &cacheOnlyTrue}
}

// Concrete merged set of configuration switches that get passed to pcidb
// internal functions
type options struct {
	chroot    string
	cacheOnly bool
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
	mergeOpts := &WithOption{}
	for _, opt := range opts {
		if opt.Chroot != nil {
			mergeOpts.Chroot = opt.Chroot
		}
		if opt.CacheOnly != nil {
			mergeOpts.CacheOnly = opt.CacheOnly
		}
	}
	// Set the default value if missing from mergeOpts
	if mergeOpts.Chroot == nil {
		mergeOpts.Chroot = &defaultChroot
	}
	if mergeOpts.CacheOnly == nil {
		mergeOpts.CacheOnly = &defaultCacheOnly
	}
	return mergeOpts
}

// New returns a pointer to a PCIDB struct which contains information you can
// use to query PCI vendor, product and class information. It accepts zero or
// more pointers to WithOption structs. If you want to modify the behaviour of
// pcidb, use one of the option modifiers when calling New. For example, to
// change the root directory that pcidb uses when discovering pciids DB files,
// call New(WithChroot("/my/root/override"))
func New(opts ...*WithOption) (*PCIDB, error) {
	mergeOpts := mergeOptions(opts...)
	useOpts := &options{
		chroot:    *mergeOpts.Chroot,
		cacheOnly: *mergeOpts.CacheOnly,
	}

	db := &PCIDB{}
	err := db.load(useOpts)
	return db, err
}
