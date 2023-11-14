package main

import (
	"flag"
	"fmt"
	"os"
)

func main () {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	tag := buildCmd.String("tag", "", "Name of the container image")
	path := buildCmd.String("path", "", "Path to Containerfile")
	buildCmd.Parse(os.Args[2:])
	build(*tag, *path)
}

func build (tag, path string) {
	fmt.Printf("Running build with tag: %s and path: %s", tag, path)
}
