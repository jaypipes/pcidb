//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package pcidb

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	PCIIDS_URI = "https://pci-ids.ucw.cz/v2.2/pci.ids.gz"
)

func (db *PCIDB) load(opts *options) error {
	cachePath := cachePath()
	// A set of filepaths we will first try to search for the pci-ids DB file
	// on the local machine. If we fail to find one, we'll try pulling the
	// latest pci-ids file from the network
	paths := []string{cachePath}
	addSearchPaths(opts, &paths)
	var foundPath string
	for _, fp := range paths {
		if _, err := os.Stat(fp); err == nil {
			foundPath = fp
			break
		}
	}

	if foundPath == "" {
		// OK, so we didn't find any host-local copy of the pci-ids DB file. Let's
		// try fetching it from the network and storing it
		if err := cacheDBFile(cachePath); err != nil {
			return err
		}
		foundPath = cachePath
	}
	f, err := os.Open(foundPath)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	return parseDBFile(db, scanner)
}

func cachePath() string {
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

// Depending on the operating system, returns a set of local filepaths to
// search for a pci.ids database file
func addSearchPaths(opts *options, paths *[]string) {
	if opts.cacheOnly {
		return
	}

	rootPath := opts.chroot

	if runtime.GOOS != "windows" {
		*paths = append(*paths, filepath.Join(rootPath, "usr", "share", "hwdata", "pci.ids"))
		*paths = append(*paths, filepath.Join(rootPath, "usr", "share", "misc", "pci.ids"))
	}
}

func ensureDir(fp string) error {
	fpDir := filepath.Dir(fp)
	if _, err := os.Stat(fpDir); os.IsNotExist(err) {
		err = os.MkdirAll(fpDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// Pulls down the latest copy of the pci-ids file from the network and stores
// it in the local host filesystem
func cacheDBFile(cacheFilePath string) error {
	ensureDir(cacheFilePath)

	response, err := http.Get(PCIIDS_URI)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	f, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	// write the gunzipped contents to our local cache file
	zr, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}
	defer zr.Close()
	if _, err = io.Copy(f, zr); err != nil {
		return err
	}
	return err
}
