//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pcidb

import (
	"github.com/jaypipes/pcidb/internal"
	"github.com/jaypipes/pcidb/types"
)

type DB = types.DB

// Backward-compat, please refer to the pcidb types.DB type definition
type PCIDB = types.DB

// New returns a pointer to a pcidb.DB struct which contains information you can
// use to query PCI vendor, product and class information.
//
// It accepts zero or more pointers to WithOption structs. If you want to
// modify the behaviour of pcidb, use one of the option modifiers when calling
// New.
//
// For example, to change the root directory that pcidb uses when discovering
// pciids DB files, call New(WithChroot("/my/root/override"))
func New(opts ...*types.WithOption) (*types.DB, error) {
	merged := internal.MergeOptions(opts...)
	scanner, err := internal.Discover(merged)
	if err != nil {
		return nil, err
	}
	return internal.FromScanner(scanner), nil
}
