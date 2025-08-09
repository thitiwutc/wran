package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func main() {
	prog := os.Args[0]
	if len(os.Args) != 2 {
		fmt.Printf("%s: invalid number of arguments\n", prog)
		printUsage(prog)

		os.Exit(1)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("%s: invalid word count\n", prog)
		printUsage(prog)
		os.Exit(1)
	}

	wordList, count := NewWordList()

	for range n {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(count)))
		if err != nil {
			fmt.Printf("%s: %s\n", prog, err)
			os.Exit(2)
		}

		head := wordList
		one := big.NewInt(1)
		for ; r.Int64() > 0; r, head = r.Sub(r, one), head.Next {
		}
		fmt.Printf("%s\n", head.Word)
	}
}

func printUsage(prog string) {
	fmt.Printf("%s: WORD_COUNT\n", prog)
}
