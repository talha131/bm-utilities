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
	"os"
	"time"

	"github.com/spf13/cobra"
)

// fileRenameCmd represents the rename command
var fileRenameCmd = &cobra.Command{
	Use:   "fileRename",
	Short: "Rename file to its ModTime",
	Long: `Rename file to the modification time of file.

Usage:

$ bmtool fileRename example.mp3 
This will rename "example.mp3" to "2016-11-04 130738.mp3"
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range args {
			fi, err := getFileInfo(e)
			if err == nil {

				n :=
					getNewName(fi.ModTime(), getFileExtension(e))
				rename(e, n)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(fileRenameCmd)
}

func rename(file string, newName string) {
	err := os.Rename(file, newName)
	if v, _ := rootCmd.Flags().GetBool("verbose"); v {
		fmt.Printf("Rename %v to %v\n", file, newName)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getNewName(t time.Time, ext string) string {
	nameFormat := "2006-01-02 150405"
	return t.Format(nameFormat) + ext
}

func getFileInfo(file string) (os.FileInfo, error) {

	var fi os.FileInfo

	// Get file stats
	fi, err := os.Stat(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return fi, err
	}

	if fi.IsDir() {
		if v, _ := rootCmd.Flags().GetBool("verbose"); v {
			fmt.Printf("Skipping %v\n", fi.Name())
		}
		return fi, errors.New("is not a file")
	}
	return fi, nil
}
