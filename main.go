package main

import (
	"fmt"
	"flag"
	"os"
)

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	if len(os.Args) <= 1 {
		fmt.Println("please enter a valid subcommand.")
		os.Exit(1)
	}

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
	case "extract":
		extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
		extractCmd.Parse(os.Args[2:])
		arg := extractCmd.Args()[0]
		extract(arg,"ubuntu-rootfs")
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}

}
