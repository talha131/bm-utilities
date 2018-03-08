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
		info, err := os.Stat(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if info.IsDir() {
			continue
		}
		time := info.ModTime().Format("2006-01-02 150405")
		newName := time + ext
		err = os.Rename(file, newName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
}
