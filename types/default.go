//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	DefaultChroot             = "/"
	DefaultCacheOnly          = false
	DefaultEnableNetworkFetch = false
)

var (
	DefaultCachePath = getCachePath()
)

func getCachePath() string {
	hdir, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed getting homedir.Dir(): %v", err)
		return ""
	}
	fp, err := homedir.Expand(filepath.Join(hdir, ".cache", "pci.ids"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed expanding local cache path: %v", err)
		return ""
	}
	return fp
}
