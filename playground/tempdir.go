package main

import (
	"os"
	"log"
	"path/filepath"
	"fmt"
)


func main() {
	dir, err := os.MkdirTemp("", "image-builder")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte("content"), 0666); err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	data, err := os.ReadFile(filepath.Join(dir, "tmpfile"))
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(data)
}
