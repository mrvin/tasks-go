// The program displays all the words that can be typed on the keyboard
// (QWERTY), moving along the adjacent keys. Neighboring keys are considered
// to have intersections with vertical and horizontal lines drawn through the
// key in question.
// For example, for D it is E, R, S, F, X, C (but not W), and for U it is Y, I,
// H, J (but not K). A word begins with any of the keys and then can only
// consist of those letters that are next to it, for example, "DESERT". Words
// from the dictionary. Words are formed by moving the keyboard. That is, if
// the beginning comes from D, then D borders on E and this is a correct
// transition, E borders on S, and so on until the transition from R to T, they
// are also neighbors, so everything converges.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var keyTable = map[rune][]rune{
	'q': {'w', 'a', 's'},
	'w': {'q', 'a', 's', 'e'},
	'e': {'w', 's', 'd', 'r'},
	'r': {'e', 'd', 'f', 't'},
	't': {'r', 'f', 'g', 'y'},
	'y': {'t', 'g', 'h', 'u'},
	'u': {'y', 'h', 'j', 'i'},
	'i': {'u', 'j', 'k', 'o'},
	'o': {'i', 'k', 'l', 'p'},
	'p': {'o', 'l'},
	'a': {'q', 'w', 's', 'z'},
	's': {'w', 'e', 'a', 'd', 'z', 'x'},
	'd': {'e', 'r', 's', 'f', 'x', 'c'},
	'f': {'r', 't', 'd', 'g', 'c', 'v'},
	'g': {'t', 'y', 'f', 'h', 'v', 'b'},
	'h': {'y', 'u', 'g', 'j', 'b', 'n'},
	'j': {'u', 'i', 'h', 'k', 'n', 'm'},
	'k': {'i', 'o', 'j', 'l', 'm'},
	'l': {'o', 'p', 'k'},
	'z': {'a', 's', 'x'},
	'x': {'z', 's', 'd', 'c'},
	'c': {'x', 'd', 'f', 'v'},
	'v': {'c', 'f', 'g', 'b'},
	'b': {'v', 'g', 'h', 'n'},
	'n': {'b', 'h', 'j', 'm'},
	'm': {'n', 'j', 'k'},
}

func main() {
	var (
		countSearchedWords uint
		maxLenWord         int
		longestWord        string
	)

	dictionaryName := flag.String("f", "dictionary.txt", "dictionary path")
	flag.Parse()
	log.Printf("Dictionary path: %s", *dictionaryName)

	dictionaryFile, err := os.Open(*dictionaryName)
	if err != nil {
		log.Fatalf("Error (os.Open) %s\n", err)
	}
	defer closeFile(dictionaryFile)

	start := time.Now()

	input := bufio.NewScanner(dictionaryFile)
	for input.Scan() {
		word := strings.ToLower(input.Text())
		if isLetter(&word) {
			if checkQwertyWord(&word) {
				fmt.Println(word)
				countSearchedWords++
				currentLen := len(word)
				if currentLen > maxLenWord {
					maxLenWord = currentLen
					longestWord = word
				}
			}
		} else {
			log.Printf("Word: %s - contains not only Latin alphabet.\n", word)
		}
	}
	if input.Err() != nil {
		log.Fatalf("Error (input.Err) %v\n", input.Err())
	}

	executionTime := time.Since(start)

	fmt.Printf("Number of searched word: %v\n", countSearchedWords)
	fmt.Printf("Max length of the searched word: %v\n", maxLenWord)
	fmt.Printf("Max length searched word: %s\n", longestWord)
	fmt.Printf("Execution time: %v s\n", executionTime.Seconds())
}

// checkQwertyWord checks can be typed on the keyboard (QWERTY) by moving along
// the adjacent keys.
func checkQwertyWord(str *string) bool {
	word := []rune(*str)
	for i, ch := range word[:len(word)-1] {
		isQwerty := false
		for _, key := range keyTable[ch] {
			if word[i+1] == key {
				isQwerty = true
				break
			}
		}
		if !isQwerty {
			return false
		}
	}

	return true
}

// isLetter checks if the word consists only of Latin lowercase letters.
func isLetter(str *string) bool {
	for _, ch := range *str {
		if ch < 'a' || ch > 'z' {
			return false
		}
	}

	return true
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalf("Error (f.Close) %v\n", err)
	}
}
