package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"path/filepath"
	"net/url"
	"net/http"
	"io"
)

func build(tag, path string) {

	distros := map[string]string{
		"ubuntu":"https://cdimage.ubuntu.com/ubuntu-base/releases/22.04/release/ubuntu-base-22.04-base-amd64.tar.gz",
	}

	file, err := os.Open(path)
	if err != nil {
		// TODO: change to err return late
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "FROM") {
			baseImage := strings.TrimSpace(strings.TrimPrefix(line, "FROM "))

			if distURL, ok := distros[baseImage]; ok {
				fmt.Println("Base Image Name:\n", distURL)
				parsedURL, err := url.Parse(distURL)
				if err != nil {
					fmt.Println("error parsing URL:", err)
					return
				}

				fileName := filepath.Base(parsedURL.Path)
				fmt.Printf("Downlading: %s \n", fileName)
				out, err := os.Create(fileName)
				if err != nil {
					fmt.Println("error creating file:", err)
					return
				}
				defer out.Close()

				resp, err := http.Get(distURL)
				if err != nil {
					fmt.Println("error making HTTP request:", err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("error: http status code %d\n", resp.StatusCode)
					return
				}

				_, err = io.Copy(out, resp.Body)
				if err != nil {
					fmt.Println("error copying to file:", err)
					return
				}

				fmt.Println("Dowload completed successfully.")

			}
		}
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
	}
	//fmt.Printf("running build with tag: %s and path: %s", tag, path)
}
