package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidRomanNum  = errors.New("invalid roman number")
	ErrInvalidArabicNum = errors.New("invalid arabic number")
)

var romanToInt = make(map[string]int)

func init() {
	romanToInt = map[string]int{
		"I":    1,
		"II":   2,
		"III":  3,
		"IV":   4,
		"V":    5,
		"VI":   6,
		"VII":  7,
		"VIII": 8,
		"IX":   9,
		"X":    10,
	}
}

func Calc(a, b int, operation string) (int, error) {
	switch operation {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	}

	return 0, ErrInvalidOperation
}

func RomanToInt(romanStrNum string) (int, error) {
	num, ok := romanToInt[romanStrNum]
	if !ok {
		return 0, ErrInvalidRomanNum
	}

	return num, nil
}

func digitToRoman(digit rune, order int) string {
	m := [4][2]rune{
		{'I', 'V'}, // 1 order
		{'X', 'L'}, // 2 order
		{'C', 'D'}, // 3 order
		{'M', ' '}, // 4 order
	}

	switch digit {
	case '1':
		return fmt.Sprintf("%c", m[order][0])
	case '2':
		return fmt.Sprintf("%[1]c%[1]c", m[order][0])
	case '3':
		return fmt.Sprintf("%[1]c%[1]c%[1]c", m[order][0])
	case '4':
		return fmt.Sprintf("%c%c", m[order][0], m[order][1])
	case '5':
		return fmt.Sprintf("%c", m[order][1])
	case '6':
		return fmt.Sprintf("%c%c", m[order][1], m[order][0])
	case '7':
		return fmt.Sprintf("%c%[2]c%[2]c", m[order][1], m[order][0])
	case '8':
		return fmt.Sprintf("%c%[2]c%[2]c%[2]c", m[order][1], m[order][0])
	case '9':
		return fmt.Sprintf("%c%c", m[order][0], m[order+1][0])
	}

	return ""
}

func IntToRoman(num int) (string, error) {
	var resultStr strings.Builder

	if num < 1 {
		return "", ErrInvalidArabicNum
	}

	numStr := strconv.Itoa(num)

	runeNumStr := []rune(numStr)
	for i, digit := range runeNumStr {
		order := len(runeNumStr) - i - 1
		resultStr.WriteString(digitToRoman(digit, order))
	}

	return resultStr.String(), nil
}
