package main

import (
	"flag"
	"fmt"
)

func main() {
	isVerbose := flag.Bool("v", false, "verbose")
	count := flag.Int("c", 3, "required count of output in seconds")
	duration := flag.Duration("d", 60, "required duration of output in seconds")
	flag.Parse()

	fmt.Println("vim-go")
}
