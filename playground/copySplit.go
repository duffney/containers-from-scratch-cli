package main

import (
	"fmt"
	"strings"
)

func main() {
	line := "COPY UBUNTU_CONTAINER_ROOT /UBUNTU_CONTAINER_ROOT"
	src := strings.Split(line, " ")[1]
	dest := strings.Split(line, " ")[2]

	fmt.Printf("parts[1]: %s \n", src)
	fmt.Printf("parts[2]: %s \n", dest)
}
