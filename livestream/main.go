package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		run(arg)
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		tag := buildCmd.String("tag", "", "Name of the container image")
		path := buildCmd.String("path", "", "Path to Containerfile")
		buildCmd.Parse(os.Args[2:])
		build(*tag, *path)
	case "child":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		child(arg)
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
	}
}

func build(tag, path string) {
	fmt.Printf("running build with tag: %s and %s\n", tag, path)
}
