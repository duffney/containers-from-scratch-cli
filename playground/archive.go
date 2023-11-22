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

func main() {
	source := "ubuntu-rootfs"
	target := "archive"

	err := TarGzip(source, target)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func TarGzip(source, target string) error {
	_, err := os.Stat(target)
	if err != nil {
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return err
		}
	}
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar.gz", filename))
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
			fmt.Println("is this the error?", err)
			return err
		})
}

func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, tarReader); err != nil {
			return err
		}
	}

	return nil
}
