package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	nameFormat := "2006-01-02 150405"
	isVerbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	for i := 0; i < len(flag.Args()); i++ {
		file := flag.Args()[i]
		ext := strings.ToLower(filepath.Ext(file))

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

		// Create new name
		time := fileInfo.ModTime().Format(nameFormat)
		newName := time + ext
		if *isVerbose {
			fmt.Printf("Rename %v to %v\n", fileInfo.Name(), newName)
		}

		// Rename
		err = os.Rename(file, newName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
}
