package main

import (
	"fmt"
	"flag"
	"os"
)

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		run(arg)
	case "child":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		child(arg)
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		tag := buildCmd.String("tag", "", "Name of container image")
		path := buildCmd.String("path", "", "Path to ContainerFile")
		buildCmd.Parse(os.Args[2:])
		build(*tag, *path)
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}

}
