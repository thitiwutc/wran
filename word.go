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
	ExactLength int
	MinLength   int
	MaxLength   int
}

func NewWordList(opts *Options) (*Node, int) {
	words := strings.Split(strings.TrimRight(wordtxt, "\n"), "\n")

	minLen := -1
	maxLen := -1
	exactLen := -1
	if opts != nil {
		minLen = opts.MinLength
		maxLen = opts.MaxLength
		exactLen = opts.ExactLength
	}

	// Nullify min length and max length filters if exact length is larger than 0.
	if exactLen > 0 {
		minLen = -1
		maxLen = -1
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

		// Filter word by maximum length
		if maxLen > 0 && len(words[i]) > maxLen {
			totalWords--
			continue
		}

		// Filter word by exact length
		if exactLen > 0 && len(words[i]) != exactLen {
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
