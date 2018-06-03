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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// videoLoopCmd represents the videoLoop command
var videoLoopCmd = &cobra.Command{
	Use:   "videoLoop",
	Short: "Concatenate same video multiple times to create a loop",
	Long: `Creates loop of a video by concatenating it multiple times.
Output format is mp4.

-c and -l are mutually exclusive. -c has precedence over -l.
`,
	Run: func(cmd *cobra.Command, args []string) {
		count, errC := cmd.Flags().GetInt("count")
		requiredLength, errD := cmd.Flags().GetInt("length")
		crossFade, _ := cmd.Flags().GetBool("withCrossFade")
		tDuration, _ := cmd.Flags().GetInt("transitionDuration")

		if errC != nil && errD != nil {
			fmt.Fprint(os.Stderr, "Unable to find Count or Length. At least one is required")
			return
		}

		if count < 2 {
			fmt.Fprint(os.Stderr, "Loop count must be at least 2.")
			return
		}

		oPath := createOutputDirectory(cmd)
		shouldConcatCountTimes := requiredLength == 0 && errC == nil && count > 2
		shouldConcatToAchieveLength := !shouldConcatCountTimes && errD == nil && requiredLength > 0

		if !crossFade {
			tDuration = 0
		}

		for _, e := range args {
			if isFileVideo(e) {
				if shouldConcatCountTimes {
					outputFileName := getOutputFileName(oPath, e, "loop", count)
					createVideoLoop(count, e, outputFileName, tDuration, crossFade)
				} else if shouldConcatToAchieveLength {
					length, err := getLength(e)
					if err != nil {
						continue
					}

					count, err := getRequiredLoopCount(length, requiredLength, tDuration)
					if err != nil {
						continue
					}

					outputFileName := getOutputFileName(oPath, e, "length", requiredLength)
					createVideoLoop(count, e, outputFileName, tDuration, crossFade)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(videoLoopCmd)

	videoLoopCmd.Flags().IntP("count", "c", 3, "Number of times to concatenate the video. Minimum 2. Default 3.")
	videoLoopCmd.Flags().IntP("length", "l", 0, "Minimum seconds of the video. Default 0 seconds")
	videoLoopCmd.Flags().BoolP("withCrossFade", "x", false, "Concatenate videos with cross fade transition. Default false.")
	videoLoopCmd.Flags().IntP("transitionDuration", "t", 2, "Transition duration. Default 2 seconds.")
	videoLoopCmd.Flags().StringP("outputDirectory", "o", "", "Output directory path. Default is current.")
}

func createVideoLoop(count int, e string, outputFileName string, tDuration int, crossFade bool) {
	if crossFade {
		createVideoLoopWithTransition(count, tDuration, e, outputFileName)
	} else {
		createVideoLoopWithoutTransition(count, e, outputFileName)
	}
}

func filterComplexWithCrossFade(count int, tDur int, length int) (filter string) {

	cf := ""
	cl := ""
	cfcl := ""
	for i := 1; i < count; i++ {
		cf = cf + fmt.Sprintf("[cf%d]", i)
		cl = cl + fmt.Sprintf("[cl%d]", i)
		cfcl = cfcl + fmt.Sprintf("[cf%d][cl%d]", i, i)
	}

	// length = 15, tDur = 5
	filter = filter + fmt.Sprintf("[0:v]trim=start=0:end=%d,setpts=PTS-STARTPTS[clip1]; ", length-tDur)               // 0 - 10
	filter = filter + fmt.Sprintf("[0:v]trim=start=%d:end=%d,setpts=PTS-STARTPTS[clip2]; ", tDur, length-tDur)        // 5 - 10
	filter = filter + fmt.Sprintf("[0:v]trim=start=%d:end=%d,setpts=PTS-STARTPTS[clip3]; ", length-tDur, length)      // 10 - 15
	filter = filter + fmt.Sprintf("[0:v]trim=start=%d:end=%d,setpts=PTS-STARTPTS[fadeoutsrc]; ", length-tDur, length) // 10 - 15
	filter = filter + fmt.Sprintf("[0:v]trim=start=0:end=%d,setpts=PTS-STARTPTS[fadeinsrc]; ", tDur)                  // 0 - 5

	filter = filter + fmt.Sprintf("[fadeinsrc]format=pix_fmts=yuva420p, fade=t=in:st=0:d=%d:alpha=1[fadein]; ", tDur)
	filter = filter + fmt.Sprintf("[fadeoutsrc]format=pix_fmts=yuva420p, fade=t=out:st=0:d=%d:alpha=1[fadeout]; ", tDur)

	filter = filter + "[fadein]fifo[fadeinfifo]; "
	filter = filter + "[fadeout]fifo[fadeoutfifo]; "
	filter = filter + "[fadeoutfifo][fadeinfifo]overlay[crossfade]; "

	filter = filter + fmt.Sprintf("[crossfade] split=%d %s ; ", count-1, cf)
	filter = filter + fmt.Sprintf("[clip2] split=%d %s ; ", count-1, cl)

	filter = filter + "[clip1]" + cfcl + "[clip3]"
	// Final number of clips to concatenate is twice of count
	filter = filter + fmt.Sprintf("concat=n=%d:v=1[output]", count*2)

	return filter
}

func createVideoLoopWithTransition(count int, tDur int, file string, outputFileName string) {
	length, err := getLength(file)
	if err != nil {
		return
	}

	if length < tDur {
		fmt.Fprint(os.Stderr, "Transition duration must be less than video length")
	}

	fc := filterComplexWithCrossFade(count, tDur, length)

	if v, _ := rootCmd.Flags().GetBool("verbose"); v {
		fmt.Printf("filter_complex is\n%s\n", fc)
	}

	cmd := exec.Command(app, "-hide_banner",
		"-i", file,
		"-f", "mp4", "-vcodec", "libx264", "-preset", "fast", "-profile:v", "main", "-movflags", "+faststart",
		"-an", "-filter_complex",
		fc,
		"-map", "[output]",
		outputFileName)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getRequiredLoopCount(length int, requiredLength int, tDuration int) (int, error) {
	if requiredLength == 0 {
		fmt.Fprintf(os.Stderr, "Required length is invalid")
		return 0, errors.New("required duration is 0")
	}

	// totalLength = count x firstClip + lastClip
	// count = (totalLength - lastClip) / firstClip
	//
	// Here totalLength is the requiredLength.
	// lastClip is transitionDuration
	// firstClip is fileLength - transitionDuration
	// count = (requiredLength - tDuration) / (length - tDuration)
	numerator := float64(requiredLength - tDuration)
	denominator := float64(length - tDuration)
	requiredLoop := int(math.Ceil(numerator / denominator))

	if v, _ := rootCmd.Flags().GetBool("verbose"); v {
		fmt.Printf("Loop %d times\n", requiredLoop)
	}

	return requiredLoop, nil
}

func getOutputFileName(oPath string, f string, t string, num int) string {
	fn := fmt.Sprintf("%s_%s-%d.mp4", getFileNameWithoutExtension(f),
		t,
		num)
	return filepath.Join(oPath, fn)
}

func createVideoLoopWithoutTransition(count int, e string, output string) {
	tmpFile, err := ioutil.TempFile(filepath.Dir(e), getFileNameWithoutExtension(e))
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpFile.Name()) // clean up

	p, _ := filepath.Abs(e)

	line := fmt.Sprintf("%s '%s'\n", "file", p)

	if v, _ := rootCmd.Flags().GetBool("verbose"); v {
		fmt.Printf("file is\n%v\n", line)
	}

	lineR := strings.Repeat(line, int(count))

	if _, err := tmpFile.WriteString(lineR); err != nil {
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	runCommandVideoLoopWithoutTransition(tmpFile.Name(),
		output)
}

func runCommandVideoLoopWithoutTransition(file string, output string) {

	cmd := exec.Command(app, "-hide_banner",
		"-f", "concat",
		"-safe", "0",
		"-i", file,
		"-qscale:v", "0",
		output)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
