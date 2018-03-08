package main

import (
	"flag"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
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

		ext := strings.ToLower(filepath.Ext(file))
		fmt.Println(mime.TypeByExtension(ext))

	}
}
