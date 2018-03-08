package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	for i := 1; i < len(os.Args); i++ {
		file := os.Args[i]
		ext := strings.ToLower(filepath.Ext(file))

		// Get file stats
		info, err := os.Stat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Skip if it is a directory
		if info.IsDir() {
			continue
		}

		// Create new name
		time := info.ModTime().Format("2006-01-02 150405")
		newName := time + ext

		// Rename
		err = os.Rename(file, newName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
}
