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
	"path/filepath"

	"github.com/spf13/cobra"
)

// audioConvertCmd represents the audioConvert command
var audioConvertCmd = &cobra.Command{
	Use:   "audioConvert",
	Short: "Convert audio file to wav or mp3",
	Long: `Convert audio file to wav or mp3 format.
Along with format conversion, it also

1. Convert stereo to mono
2. Set audio sample frequency to 44100
3. Set mp3 bit rate to 32k

It creates output in the -o directory with same name except the new extension.
If -o is not given then it creates output in the same directory.

It removes album art from the audio. It only picks the audio stream.

Usage:
$ bmtool audio convert -f mp3 -o eg example.wav 

It will convert "example.wav" to "example.mp3" in ./eg directory
`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		if format != "mp3" && format != "wav" {
			fmt.Fprintf(os.Stderr, "Unknown format %v. Valid values are [mp3|wav]\n", format)
			return
		}

		oPath := createOutputDirectory(cmd)

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

				fn := getFileNameWithoutExtension(e) + "." + format

				output = append(output, "-map",
					fmt.Sprintf("%d:a", len(input)/3), filepath.Join(oPath, fn))
			}
		}

		if len(input) > 0 && len(output) > 0 {
			convertFile(input, output)
		}

	},
}

func init() {
	rootCmd.AddCommand(audioConvertCmd)

	audioConvertCmd.Flags().StringP("format", "f", "wav", "Output format. [mp3|wav]")
	audioConvertCmd.Flags().StringP("outputDirectory", "o", "", "Output directory path. Default is current.")
}

func convertFile(input []string, output []string) {

	var a []string

	a = append(a, "-hide_banner")
	a = append(a, input...)

	a = append(a, output...)

	if v, _ := rootCmd.Flags().GetBool("verbose"); v {
		fmt.Printf("Command is\n%s %v\n", app, a)
	}

	cmd := exec.Command(app, a...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
