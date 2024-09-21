//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

var (
	trueVar = true
)

// WithOption is used to represent optionally-configured settings
type WithOption struct {
	// Chroot is the directory that pcidb uses when attempting to discover
	// pciids DB files
	Chroot *string
	// CacheOnly is mostly just useful for testing. It essentially disables
	// looking for any non ~/.cache/pci.ids filepaths (which is useful when we
	// want to test the fetch-from-network code paths
	CacheOnly *bool
	// CachePath overrides the pcidb cache path, which defaults to
	// $HOME/.cache/pci.ids
	CachePath *string
	// Enables fetching a pci-ids from a known location on the network if no
	// local pci-ids DB files can be found.
	EnableNetworkFetch *bool
	// Path points to the absolute path of a pci.ids file in a non-standard
	// location.
	Path *string
}

func WithChroot(dir string) *WithOption {
	return &WithOption{Chroot: &dir}
}

func WithCachePath(path string) *WithOption {
	return &WithOption{CachePath: &path}
}

func WithCacheOnly() *WithOption {
	return &WithOption{CacheOnly: &trueVar}
}

func WithDirectPath(path string) *WithOption {
	return &WithOption{Path: &path}
}

func WithEnableNetworkFetch() *WithOption {
	return &WithOption{EnableNetworkFetch: &trueVar}
}
