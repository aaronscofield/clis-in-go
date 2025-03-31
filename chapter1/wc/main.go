package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Calling the count function to count the number of words
	// received from the standard inmput and printing it out

	// Defining a boolean flag -l to count lines instead of words
	// flag.Bool returns a pointer to a boolean
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	flag.Parse()

	// Ensure XOR logic: either lines or bytes must be set, but not both
	if *lines && *bytes {
		fmt.Fprintln(os.Stderr, "Error: You must specify either -l (lines) or -b (bytes), but not both.")
		os.Exit(1)
	}

	fmt.Println(count(os.Stdin, *lines, *bytes)) // Destructure the boolean pointer
}

func count(r io.Reader, countLines bool, countBytes bool) int {
	// A scanner is used to read text from a Reader (files, etc)
	scanner := bufio.NewScanner(r)

	// If the countlines flag is not set, we want to count words
	// the default scanner is split by lines
	if countBytes {
		scanner.Split(bufio.ScanBytes)
	} else if countLines {
		scanner.Split(bufio.ScanLines)
	} else {
		scanner.Split(bufio.ScanWords)
	}

	// Defining a counter
	wc := 0

	// For every word scanned, increment the counter
	for scanner.Scan() {
		wc++
	}

	// Return the total
	return wc
}
