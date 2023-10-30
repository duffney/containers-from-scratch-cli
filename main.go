
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
		runContainer(arg)
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}

}

func runContainer(args []string) {
	fmt.Printf("Running %v \n", args)
}
