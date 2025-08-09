package main

import (
	_ "embed"
	"strings"
)

//go:embed words.txt
var wordtxt string

type Node struct {
	Word string
	Next *Node
}

type Options struct {
	MinLength int
}

func NewWordList(opts *Options) (*Node, int) {
	words := strings.Split(strings.TrimRight(wordtxt, "\n"), "\n")

	minLen := -1
	if opts != nil {
		minLen = opts.MinLength
	}

	var head *Node
	var prev *Node
	totalWords := len(words)

	for i := 0; i < len(words); i++ {
		// Filter word by minimum length
		if minLen > 0 && len(words[i]) < minLen {
			totalWords--
			continue
		}

		cur := Node{
			Word: words[i],
			Next: nil,
		}

		if head == nil {
			head = &cur
			prev = head
			continue
		}

		prev.Next = &cur
		prev = &cur
	}

	return head, totalWords
}
