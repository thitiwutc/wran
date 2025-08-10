package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
)

func main() {
	isDup := flag.Bool("dup", false, "Allow duplicate words if true")
	minLen := flag.Int("minlen", -1, "Minimum word length. minlen < 0 allows any lengths")
	maxLen := flag.Int("maxlen", -1, "Maximum word length. maxlen < 0 allows any lengths")
	help := flag.Bool("h", false, "Print help message")
	printVer := flag.Bool("v", false, "Print version")
	flag.Parse()

	prog := os.Args[0]

	if *help {
		printUsage(prog, os.Stdout)
		os.Exit(0)
	}

	if *printVer {
		fmt.Println(version())
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("%s: invalid number of arguments\n", prog)
		printUsage(prog, os.Stderr)

		os.Exit(1)
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("%s: invalid word count\n", prog)
		printUsage(prog, os.Stderr)
		os.Exit(1)
	}

	wordList, count := NewWordList(&Options{
		MinLength: *minLen,
		MaxLength: *maxLen,
	})
	if count == 0 {
		fmt.Printf("%s: no words after filter", prog)
		os.Exit(2)
	}

	// Random n word(s)
	for range n {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(count)))
		if err != nil {
			fmt.Printf("%s: %s\n", prog, err)
			os.Exit(3)
		}

		cur := wordList
		one := big.NewInt(1)
		var prev *Node
		for ; r.Int64() > 0; r, cur, prev = r.Sub(r, one), cur.Next, cur {
		}

		// Remove the random word from the list to prevent word duplication.
		if !*isDup {
			if prev == nil {
				// Remove the first word.
				wordList = wordList.Next
			} else if cur.Next == nil {
				// Remove the last word.
				prev.Next = nil
			} else {
				// Remove the middle word.
				prev.Next = cur.Next
			}
			count--

			// No more word in list.
			if count == 0 {
				break
			}
		}

		fmt.Printf("%s\n", cur.Word)
	}
}

func version() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	matched, _ := regexp.MatchString(`v\d+\.\d+\.\d+`, buildInfo.Main.Version)
	if matched {
		return buildInfo.Main.Version
	}

	return "unknown"
}

func printUsage(prog string, w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [OPTION]... WORD_COUNT\n\n", prog)
	fmt.Fprintf(w, "Options:\n")
	flag.PrintDefaults()
}
