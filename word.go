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

func NewWordList() (*Node, int) {
	words := strings.Split(strings.TrimRight(wordtxt, "\n"), "\n")

	head := new(Node)
	head.Word = words[0]
	prev := head

	for i := 1; i < len(words); i++ {
		cur := Node{
			Word: words[i],
			Next: nil,
		}

		prev.Next = &cur
		prev = &cur
	}

	return head, len(words)
}
