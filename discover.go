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
)

const (
	PCIIDS_URI = "https://pci-ids.ucw.cz/v2.2/pci.ids.gz"
)

func (db *PCIDB) load(ctx *context) error {
	var foundPath string
	var foundCompressedPath string
	for _, fp := range ctx.searchPaths {
		if _, err := os.Stat(fp); err == nil {
			foundPath = fp
			break
		}

		if _, err := os.Stat(fp + ".gz"); err == nil {
			foundCompressedPath = fp
			break
		}
	}

	if foundPath == "" {
		// OK, so we didn't find any host-local copy of the pci-ids DB file. If
		// we found a local compressed copy, we'll use that. If not, we fetch the
		// latest from the network.
		if err := cacheDBFile(ctx.cachePath, foundCompressedPath, ctx.localOnly); err != nil {
			return err
		}
		foundPath = ctx.cachePath
	}

	f, err := os.Open(foundPath)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	return parseDBFile(db, scanner)
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
func cacheDBFile(cacheFilePath string, compressedFilePath string, localOnly bool) error {
	ensureDir(cacheFilePath)
	var response *http.Response
	var err error

	if localOnly {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c := &http.Client{Transport: t}
		if compressedFilePath == "" {
			return fmt.Errorf("failed to locate compressed pci-ids.")
		}
		response, err = c.Get("file://" + compressedFilePath)
	} else {
		response, err = http.Get(PCIIDS_URI)
		if err != nil {
			return err
		}
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
