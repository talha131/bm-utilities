package main

import (
	"flag"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	isVerbose := flag.Bool("v", false, "verbose")
	outputFormat := flag.String("o", "wav", "output format. wav|mp3")
	flag.Parse()

	cmd := "ffmpeg"
	cmdArgs := "-ac 1 -ab 64k -ar 44100"

	for i := 0; i < len(flag.Args()); i++ {
		file := flag.Args()[i]

		if isFileAudio(file, *isVerbose) {
			name := strings.TrimSuffix(file, filepath.Ext(file))
			output := name + "." + *outputFormat
			_, err := exec.Command(cmd, "-i", file, cmdArgs, output).Output()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
	}
}

func isFileAudio(file string, isVerbose bool) bool {
	mimeType := "audio/mpeg"
	// Get file stats
	fileInfo, err := os.Stat(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	if fileInfo.IsDir() {
		if isVerbose {
			fmt.Printf("Skipping %v\n", fileInfo.Name())
		}
		return false
	}

	ext := strings.ToLower(filepath.Ext(file))
	fileType := mime.TypeByExtension(ext)
	if fileType != mimeType {
		if isVerbose {
			fmt.Printf("Skipping %v \t %v\n", fileInfo.Name(), fileType)
		}
		return false
	}

	return true
}
