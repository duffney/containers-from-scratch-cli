package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func createImage(source, target, tag string) error {
	_, err := os.Stat(target)
	if err != nil {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return err
		}
	}
	//filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar.gz", tag))
	tarGzFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarGzFile.Close()

	gzWriter := gzip.NewWriter(tarGzFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	_, err = os.Stat(source)
	if err != nil {
		return nil
	}

	processedSymlinks := make(map[string]bool)

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip the base directory itself
			if path == source {
				return nil
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			header.Name = filepath.ToSlash(strings.TrimPrefix(path, source+"/"))

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			// skip directories
			if info.IsDir() {
				return nil
			}

			// store symlink as symlink
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				link, err := os.Readlink(path)
				if err != nil {
					return fmt.Errorf("failed to read symlink %s: %v", path, err)
				}
				
				if processedSymlinks[link] {
					return nil
				}

				processedSymlinks[link] = true
				header.Linkname = link

				if err := tarWriter.WriteHeader(header); err != nil {
					return err
				}

				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarWriter, file)
			return err
		})
}
