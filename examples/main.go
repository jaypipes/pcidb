package main

import (
	"cmp"
	"fmt"
	"slices"
	"sort"

	"github.com/jaypipes/pcidb"
)

type ByCountSeparateSubvendors []*pcidb.Product

func (v ByCountSeparateSubvendors) Len() int {
	return len(v)
}

func (v ByCountSeparateSubvendors) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v ByCountSeparateSubvendors) Less(i, j int) bool {
	iVendor := v[i].VendorID
	iSetSubvendors := make(map[string]bool, 0)
	iNumDiffSubvendors := 0
	jVendor := v[j].VendorID
	jSetSubvendors := make(map[string]bool, 0)
	jNumDiffSubvendors := 0

	for _, sub := range v[i].Subsystems {
		if sub.VendorID != iVendor {
			iSetSubvendors[sub.VendorID] = true
		}
	}
	iNumDiffSubvendors = len(iSetSubvendors)

	for _, sub := range v[j].Subsystems {
		if sub.VendorID != jVendor {
			jSetSubvendors[sub.VendorID] = true
		}
	}
	jNumDiffSubvendors = len(jSetSubvendors)

	return iNumDiffSubvendors > jNumDiffSubvendors
}

func main() {
	pci, err := pcidb.New()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	for _, devClass := range pci.Classes {
		fmt.Printf(" Device class: %v ('%v')\n", devClass.Name, devClass.ID)
		for _, devSubclass := range devClass.Subclasses {
			fmt.Printf("    Device subclass: %v ('%v')\n", devSubclass.Name, devSubclass.ID)
			for _, progIface := range devSubclass.ProgrammingInterfaces {
				fmt.Printf("        Programming interface: %v ('%v')\n", progIface.Name, progIface.ID)
			}
		}
	}

	vendors := make([]*pcidb.Vendor, len(pci.Vendors))
	x := 0
	for _, vendor := range pci.Vendors {
		vendors[x] = vendor
		x++
	}

	slices.SortFunc(vendors, func(a, b *pcidb.Vendor) int {
		return cmp.Compare(len(a.Products), len(b.Products))
	})
	slices.Reverse(vendors)

	fmt.Println("Top 5 vendors by product")
	fmt.Println("====================================================")
	for _, vendor := range vendors[0:5] {
		fmt.Printf("%v ('%v') has %d products\n", vendor.Name, vendor.ID, len(vendor.Products))
	}

	products := make([]*pcidb.Product, len(pci.Products))
	x = 0
	for _, product := range pci.Products {
		products[x] = product
		x++
	}

	sort.Sort(ByCountSeparateSubvendors(products))

	fmt.Println("Top 2 products by # different subvendors")
	fmt.Println("====================================================")
	for _, product := range products[0:2] {
		vendorID := product.VendorID
		vendor := pci.Vendors[vendorID]
		setSubvendors := make(map[string]bool, 0)

		for _, sub := range product.Subsystems {
			if sub.VendorID != vendorID {
				setSubvendors[sub.VendorID] = true
			}
		}
		fmt.Printf("%v ('%v') from %v\n", product.Name, product.ID, vendor.Name)
		fmt.Printf(" -> %d subsystems under the following different vendors:\n", len(setSubvendors))
		for subvendorID, _ := range setSubvendors {
			subvendor, exists := pci.Vendors[subvendorID]
			subvendorName := "Unknown subvendor"
			if exists {
				subvendorName = subvendor.Name
			}
			fmt.Printf("      - %v ('%v')\n", subvendorName, subvendorID)
		}
	}
}
