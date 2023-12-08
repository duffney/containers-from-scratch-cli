package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mholt/archiver/v4"
)

func main() {
	file, err := os.Open("example.tar.gz")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	handler := func(ctx context.Context, f archiver.File) error {
		return nil
	}

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}

	fileList := []string{"/bin", ""}

	err = format.Extract(context.Background(), file, fileList, handler)
	if err != nil {
		fmt.Println(err)
	}
}
