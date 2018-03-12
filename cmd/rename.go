// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename file to its ModTime",
	Long: `Rename file to the modification time of file.

	Example:

	$ bmtool file rename example.mp3 
	This will rename "example.mp3" to "2016-11-04 130738.mp3"
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range args {
			fi, err := getFileInfo(e)
			if err == nil {

				ext := strings.ToLower(filepath.Ext(e))
				n :=
					getNewName(fi.ModTime(), ext)
				rename(e, n)
			}

		}
	},
}

func init() {
	fileCmd.AddCommand(renameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func rename(file string, newName string) {
	err := os.Rename(file, newName)
	if Verbose {
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
		if Verbose {
			fmt.Printf("Skipping %v\n", fi.Name())
		}
		return fi, errors.New("Is not a file")
	}
	return fi, nil
}
