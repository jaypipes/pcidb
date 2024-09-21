//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package internal

import (
	"testing"

	"github.com/jaypipes/pcidb/types"
)

func TestMergeOptions(t *testing.T) {
	// Verify the default values are set if no overrides are passed
	opts := MergeOptions()
	if opts.Chroot == nil {
		t.Fatalf("Expected opts.Chroot to be non-nil.")
	}
	if opts.CacheOnly == nil {
		t.Fatalf("Expected opts.CacheOnly to be non-nil.")
	}
	if opts.EnableNetworkFetch == nil {
		t.Fatalf("Expected opts.EnableNetworkFetch to be non-nil.")
	}
	if opts.Path == nil {
		t.Fatalf("Expected opts.DirectPath to be non-nil.")
	}

	// Verify if we pass an override, that value is used not the default
	opts = MergeOptions(types.WithChroot("/override"))
	if opts.Chroot == nil {
		t.Fatalf("Expected opts.Chroot to be non-nil.")
	} else if *opts.Chroot != "/override" {
		t.Fatalf("Expected opts.Chroot to be /override.")
	}

	opts = MergeOptions(types.WithDirectPath("/mnt/direct/pci.ids"))
	if opts.Path == nil {
		t.Fatalf("Expected opts.DirectPath to be non-nil.")
	} else if *opts.Path != "/mnt/direct/pci.ids" {
		t.Fatalf("Expected opts.DirectPath to be /mnt/direct/pci.ids")
	}
}
