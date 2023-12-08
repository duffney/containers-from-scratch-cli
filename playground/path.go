package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	//tempDir := "/tmp/dir"
	curDir := "/home/jduffney/github/containers-from-scratch-cli"
	tag := "container"
	
	imagePath := filepath.Join(curDir,tag+".tar.gz")	
	fmt.Println(imagePath)
}
