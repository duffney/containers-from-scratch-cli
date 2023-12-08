/*
Download distro
Create temp dir
Extract rootfs
Run COPY instructions
Create container archive
(Let run func read RUN instruction)
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"github.com/mholt/archiver/v3"
)

func build(tag, path string) {

	distros := map[string]string{
		"ubuntu": "https://cdimage.ubuntu.com/ubuntu-base/releases/22.04/release/ubuntu-base-22.04-base-amd64.tar.gz",
	}

	file, err := os.Open(path)
	if err != nil {
		// TODO: change to err return late
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//asign current dir to var
	curDir := filepath.Dir(path)

	// Create temp dir
	tempDir, err := os.MkdirTemp("", "image-builder")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "FROM") {
			baseImage := strings.TrimSpace(strings.TrimPrefix(line, "FROM "))

			//download distro
			if distURL, ok := distros[baseImage]; ok {
				fmt.Println("Base Image Name:\n", distURL)
				parsedURL, err := url.Parse(distURL)
				if err != nil {
					fmt.Println("error parsing URL:", err)
					return
				}

				fileName := filepath.Base(parsedURL.Path)
				filePath := filepath.Join(tempDir, fileName)

				fmt.Printf("Downlading to: %s \n", filePath)
				out, err := os.Create(filePath)
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

				// TODO: extract rootfs from distro archive
				err = archiver.Unarchive(filePath, tempDir)
				if err != nil {
					fmt.Println("error unarchiving distro:", err)
					return
				}


				// remove distro archive, prevent from including in container archive
				err = os.Remove(filePath)
				if err != nil {
					fmt.Println("error removing file:", err)
					return
				}

			} else {
				fmt.Println("distro not supported")
			}
		}

		if strings.HasPrefix(line, "COPY") {
			src := strings.Split(line, " ")[1]
			dest := strings.Split(line, " ")[2]

			srcFile, err := os.Open(filepath.Join(curDir, src))
			if err != nil {
				fmt.Println("error opening file:", err)
				return
			}
			defer srcFile.Close()

			destFile, err := os.Create(filepath.Join(tempDir, dest))
			if err != nil {
				fmt.Println("error creating file:", err)
				return
			}

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				fmt.Println("error copying file:", err)
				return
			}

			fmt.Printf("%s: \n", line)
			fmt.Printf("COPY instruction found")
		}
		fmt.Println(line)

		// if line == "RUN"
		// TODO: run command
		// how could I pass this or store it for the run cmd? config file?
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
	}

	// create container archive
	//createImage(tempDir, curDir, tag)
	imagePath := filepath.Join(curDir,tag+".tar.gz")	
	err = archiver.Archive([]string{tempDir+"/"},imagePath)
	if err != nil {
		fmt.Println("error archiving image:", err)
		return
	}
}
