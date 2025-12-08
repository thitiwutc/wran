package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const maxWordsPerLine = 10

func main() {
	isDup := flag.Bool("dup", false, "Allow duplicate words if true")
	minLen := flag.Int("minlen", -1, "Minimum word length. minlen < 0 allows any lengths")
	maxLen := flag.Int("maxlen", -1, "Maximum word length. maxlen < 0 allows any lengths")
	help := flag.Bool("h", false, "Print help message")
	printVer := flag.Bool("V", false, "Print version")
	flag.Parse()

	prog := "wran"

	if *help {
		printUsage(prog, os.Stdout)
		os.Exit(0)
	}

	if *printVer {
		fmt.Println(version())
		os.Exit(0)
	}

	wordCount := 10
	args := flag.Args()
	// Use default word count.
	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("%s: invalid word count\n", prog)
			printUsage(prog, os.Stderr)
			os.Exit(1)
		}

		wordCount = n
	} else if len(args) > 1 {
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

	randWords := make([]string, 0, wordCount)

	// Random n word(s)
	for range wordCount {
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

		randWords = append(randWords, cur.Word)
	}

	var longestWordLen int
	for _, word := range randWords {
		longestWordLen = max(longestWordLen, len(word))
	}
	sort.Strings(randWords)

	var sb strings.Builder

	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		width, _, err := term.GetSize(fd)
		if err != nil {
			fmt.Printf("get terminal size failed: %v\n", err)
			os.Exit(4)
		}

		itemsPerLine := float64(width) / float64(longestWordLen) * 0.6
		itemsPerLine = min(itemsPerLine, maxWordsPerLine)
		itemsPerLineInt := int(math.Floor(itemsPerLine))

		for i, word := range randWords {
			if i > 0 {
				if i%itemsPerLineInt == 0 {
					sb.WriteRune('\n')
				} else {
					sb.Write([]byte("  "))
				}
			}

			if wordCount > itemsPerLineInt {
				sb.WriteString(fmt.Sprintf("%-*s", longestWordLen, word))
			} else {
				// No padding if output is only 1 line.
				sb.WriteString(word)
			}
		}

		fmt.Println(sb.String())
	} else {
		// Print newline-separated words for non-terminal output.
		for _, word := range randWords {
			sb.WriteString(word)
			sb.WriteRune('\n')
		}

		fmt.Print(sb.String())
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
	fmt.Fprintf(w, "Usage: %s [OPTIONS] [WORD_COUNT]\n\n", prog)
	fmt.Fprintf(w, "Options:\n")
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
}
