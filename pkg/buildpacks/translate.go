package buildpacks

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cnbapplifecycle/pkg/archive"
	"code.cloudfoundry.org/cnbapplifecycle/pkg/log"
	"github.com/cespare/xxhash/v2"
)

func Translate(bps []string, buildpacksDir string, logger *log.Logger) ([]string, error) {
	newList := []string{}

	for _, bp := range bps {
		bpDir := buildpackPath(bp, buildpacksDir)
		downloaded, err := checkIfDownloaded(bpDir)
		if err != nil {
			return nil, err
		}

		if downloaded {
			newPath, err := createArchive(bpDir)
			if err != nil {
				return nil, err
			}

			newList = append(newList, newPath)
		} else {
			newList = append(newList, bp)
		}
	}

	return newList, nil
}

func checkIfDownloaded(path string) (bool, error) {
	fi, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	if !fi.IsDir() {
		return false, fmt.Errorf("%s is not a directory", path)
	}

	return true, nil
}

func buildpackPath(name string, path string) string {
	return filepath.Join(path, fmt.Sprintf("%016x", xxhash.Sum64String(name)))
}

func createArchive(path string) (string, error) {
	newPath := fmt.Sprintf("%s.tgz", path)

	f, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bgw := gzip.NewWriter(f)
	defer bgw.Close()

	return fmt.Sprintf("file://%s", newPath), archive.FromDirectory(path, tar.NewWriter(bgw))
}
