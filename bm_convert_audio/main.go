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

	for i := 0; i < len(flag.Args()); i++ {
		file := flag.Args()[i]

		if isFileAudio(file, *isVerbose) {

			name := strings.TrimSuffix(file, filepath.Ext(file))
			output := name + "." + *outputFormat
			convertFile(file, output)
		}
	}
}

func convertFile(input string, output string) {

	cmd := exec.Command("ffmpeg",
		"-i", input,
		"-ac", "1",
		"-ab", "64k",
		"-ar", "44100",
		output)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func isFileAudio(file string, isVerbose bool) bool {
	mimeTypeMp3 := "audio/mpeg"
	mimeTypeWav := "audio/x-wav"
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
	if fileType != mimeTypeMp3 && fileType != mimeTypeWav {
		if isVerbose {
			fmt.Printf("Skipping %v \t %v\n", fileInfo.Name(), fileType)
		}
		return false
	}

	return true
}
