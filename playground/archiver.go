package main

import (
	"context"
	"os"
	"fmt"
	"github.com/mholt/archiver/v4"
)

func main() {
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		"ubuntu-rootfs/": "",
	})

	if err != nil {
		fmt.Println(err)
	}

	out, err := os.Create("example.tar.gz")
	if err != nil {
		fmt.Println(err)
	}

	defer out.Close()

	format := archiver.CompressedArchive {
		Compression: archiver.Gz{},
		Archival: archiver.Tar{},
	}

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		fmt.Println(err)
	} 
	
	handler := func(ctx context.Context, f archiver.File) error {
		return nil
	}
	
	format.Extract(context.Background(), out, nil, handler)
	if err != nil {
		fmt.Println(err)
	}
}

