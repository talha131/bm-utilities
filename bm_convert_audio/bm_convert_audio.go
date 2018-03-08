package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {

	isVerbose := flag.Bool("v", false, "verbose")
	// outputFormat := flag.String("o", "wav", "output format. wav|mp3")
	flag.Parse()

	for i := 0; i < len(flag.Args()); i++ {
		file := flag.Args()[i]

		// Get file stats
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Skip if it is a directory
		if fileInfo.IsDir() {
			if *isVerbose {
				fmt.Printf("Skipping %v\n", fileInfo.Name())
			}
			continue
		}

		fmt.Println(getFileContentType(file))

	}
}

func getFileContentType(fileName string) (string, error) {

	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
