// Copyright © 2018 Talha Mansoor <talha131@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mimeTypeMp3 = "audio/mpeg"
	mimeTypeWav = "audio/x-wav"
	app         = "ffmpeg"
	wavOption   = []string{"-ac", "1", "-ar", "44100"}
	mp3Option   = []string{"-ac", "1", "-ar", "44100", "-b:a", "32k"}
)

// getFileExtension returns file extension from file name
func getFileExtension(file string) string {
	return strings.ToLower(filepath.Ext(file))
}

// getFileNameWithoutExtension returns file name sans extension
func getFileNameWithoutExtension(file string) string {
	file = filepath.Base(file)
	return strings.TrimSuffix(file, filepath.Ext(file))
}

// isFileAudio checks if file is mp3 or wav using mime type
func isFileAudio(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	if fi.IsDir() {
		if v, _ := rootCmd.Flags().GetBool("verbose"); v {
			fmt.Printf("%v \t direcotry\n", fi.Name())
		}
		return false
	}

	fileType := mime.TypeByExtension(getFileExtension(file))

	if fileType != mimeTypeMp3 && fileType != mimeTypeWav {
		if v, _ := rootCmd.Flags().GetBool("verbose"); v {
			fmt.Printf("%v \t %v\n", fi.Name(), fileType)
		}
		return false
	}

	return true
}

func createDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func createOutputDirectory(cmd *cobra.Command) string {
	o, _ := cmd.Flags().GetString("outputDirectory")
	if o != "" {
		createDirectory(o)
	}
	return o
}
