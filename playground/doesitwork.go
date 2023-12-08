package main

import (
	"github.com/mholt/archiver/v3"
)

func main() {
	err := archiver.Archive([]string{"ubuntu-rootfs/"}, "example.tar.gz")
	if err != nil {
		panic(err)
	}

	err = archiver.Unarchive("example.tar.gz", "example/")
	if err != nil {
		panic(err)
	}

}
