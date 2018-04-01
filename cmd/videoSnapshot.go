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

// videoSnapshotCmd represents the videoSnapshot command
var videoSnapshotCmd = &cobra.Command{
	Use:   "videoSnapshot",
	Short: "Takes snapshot of video at 2nd second",
	Long: `Takes snapshot of video.
Default is to take snapshot at 2nd second i.e. 00::00::2.0

if -m flag is used then snapshot is taken from the middle of the video.
If video is 30 minute long, then it will take snaptshot right at 00:15:00.

Default output format is png. If output format is set to jpeg then it is exported
at highest quality.
`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		if format != "png" && format != "jpg" {
			fmt.Fprintf(os.Stderr, "Unknown format %v. Valid values are [png|jpg]\n", format)
			return
		}

		oPath := createOutputDirectory(cmd)
		for _, e := range args {
			if isFileVideo(e) {
				timestamp := "2"
				if d, e := getLength(e); e != nil || d < 2 {
					continue
				}

				if m, _ := cmd.Flags().GetBool("mid"); m {

					var err error
					timestamp, err = getMidTimestamp(e)
					if err != nil {
						continue
					}
				}

				f := fmt.Sprintf("%s-%s.%s", getFileNameWithoutExtension(e), timestamp, format)
				of := filepath.Join(oPath, f)
				createVideoSnapshot(timestamp, e, of)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(videoSnapshotCmd)
	videoSnapshotCmd.Flags().BoolP("mid", "m", false, "Take snapshot from mid")
	videoSnapshotCmd.Flags().StringP("format", "f", "png", "Output format. png|jpg")
	videoSnapshotCmd.Flags().StringP("outputDirectory", "o", "", "Output directory path. Default is current.")
}

func getMidTimestamp(file string) (string, error) {
	fileDuration, err := getLength(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get duration of %s\t%s", file, err)
		return "00", err
	}

	mid := fmt.Sprintf("%d", fileDuration/2)

	if v, _ := rootCmd.Flags().GetBool("verbose"); v {
		fmt.Printf("Total %d, mid %s", fileDuration, mid)
	}

	return mid, nil
}

func createVideoSnapshot(timestamp string, file string, output string) {
	cmd := exec.Command(app, "-hide_banner",
		"-ss", timestamp,
		"-i", file,
		"-vframes", "1",
		"-qscale:v", "1", // meaningless if output format is png
		output)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
