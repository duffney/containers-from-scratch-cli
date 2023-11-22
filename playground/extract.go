package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	src := "ubuntu-base-22.04-base-amd64.tar.gz"
	dest := "ubuntu-rootfs"

	err := extract(src, dest)
	if err != nil {
		log.Fatal("Extraction failed:", err)
	}
}

func extract(src, dest string) error {
	gzipFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Failed to open gzip file: %v", err)
	}
	defer gzipFile.Close()

	_, err = os.Stat(dest)
	if err != nil {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return fmt.Errorf("Failed to create extract directory: %v", err)
		}
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return fmt.Errorf("Failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("Failed to read header: %v", err)
		}

		extractPath := fmt.Sprintf("%s/%s", dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir: // if a dir, create it
			err := os.MkdirAll(extractPath, 0755)
			if err != nil {
				return fmt.Errorf("Failed to create extract directory: %v", err)
			}
		case tar.TypeReg: // if a file, create it and write its contents
			extractFile, err := os.OpenFile(extractPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("Failed to create extract file: %v", err)
			}
			_, err = io.Copy(extractFile, tarReader)
			extractFile.Close()
			if err != nil {
				return fmt.Errorf("Error copying file contents: %v", err)
			}
		case tar.TypeSymlink: // if a symlink, create it
			err := os.Symlink(header.Linkname, extractPath)
			if err != nil {
				return fmt.Errorf("Failed to create symbolic link: %v", err)
			}
		}
	}

	return nil
}
