package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mholt/archiver/v4"
)

func main() {
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		"fileToArchive": "",
		"cmds.go":       "",
	})

	if err != nil {
		fmt.Println(err)
	}

	out, err := os.Create("example.tar.gz")
	if err != nil {
		fmt.Println(err)
	}

	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Open("example.tar.gz")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	handler := func(ctx context.Context, f archiver.File) error {
		fmt.Println(f.Name())
		return nil
	}

	err = format.Extract(context.Background(), file, nil, handler)
	if err != nil {
		fmt.Println(err)
	}
}
