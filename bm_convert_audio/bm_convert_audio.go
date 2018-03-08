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

	var input []string
	var output []string

	for i := 0; i < len(flag.Args()); i++ {
		file := flag.Args()[i]

		if isFileAudio(file, *isVerbose) {
			input = append(input, "-i")
			input = append(input, file)

			name := strings.TrimSuffix(file, filepath.Ext(file))
			outputName := name + "." + *outputFormat

			output = append(output, "-map")
			output = append(output, fmt.Sprintf("%v", len(output)/4))

			output = append(output, "-f")
			output = append(output, *outputFormat)

			output = append(output, outputName)

		}
	}

	convertFiles(input, output)
}

func convertFiles(input []string, output []string) {

	cmd := exec.Command("ffmpeg",
		"-ac", "1",
		"-ab", "64k",
		"-ar", "44100",
		strings.Join(input, " "),
		strings.Join(output, " "))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
