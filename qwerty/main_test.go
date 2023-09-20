package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type pairWordResult struct {
	word   string
	result bool
}

func TestIsLetter(t *testing.T) {
	var words = []pairWordResult{
		{"123445", false},
		{"heLl", false},
		{"it's", false},
		{"boredom", true},
		{"hello", true},
		{"red", true},
	}

	for _, pair := range words {
		require.Equal(t, isLetter(&pair.word), pair.result, "Checking Latin lowercase letters: %s", pair.word)
	}
}

func TestCheckQwertyWord(t *testing.T) {
	var words = []pairWordResult{
		{"qwertyuioplzjhgfdsazxcvbnm", false},
		{"qwertyi", false},
		{"qrfvbgtyhn", false},
		{"qazxswedcvfrtgbnhyujmkiolp", true},
		{"qwerty", true},
		{"hello", false},
	}

	for _, pair := range words {
		require.Equal(t, checkQwertyWord(&pair.word), pair.result, "Checking Latin lowercase letters: %s", pair.word)
	}

}
