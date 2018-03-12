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
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert audio file to wav or mp3",
	Long: `Convert audio file to wav or mp3 format.
Alongwith format conversion, it also

1. Convert stereo to mono
2. Set audio sample frequency to 44100
3. Set mp3 bit rate to 64k

It creates output in the same directory with same name except the new extension.
You must make sure directory does not already have a file with the same name.

Usage:
$ bmtool audio convert -f mp3 example.wav

It will convert "example.wav" to "example.mp3"
`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		if format != "mp3" && format != "wav" {
			fmt.Fprintf(os.Stderr, "Unknown format %v. Valid values are [mp3|wav]\n", format)
			return
		}

		var (
			input  []string
			output []string
		)

		for _, e := range args {
			if isFileAudio(e) {
				input = append(input, "-i", e)

				if format == "wav" {
					output = append(output, wavOption...)
				} else {
					output = append(output, mp3Option...)
				}

				output = append(output, "-map",
					fmt.Sprintf("%d", len(input)/3),
					getFileNameWithoutExtension(e)+"."+format)
			}
		}

		if len(input) > 0 && len(output) > 0 {
			convertFile(input, output)
		}

	},
}

func init() {
	audioCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringP("format", "f", "wav", "Output format. [mp3|wav]")
}

func convertFile(input []string, output []string) {

	var a []string

	a = append(a, input...)

	a = append(a, output...)

	cmd := exec.Command(app, a...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
