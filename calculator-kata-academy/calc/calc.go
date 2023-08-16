package calc

import "errors"

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidRomanNum  = errors.New("invalid roman number")
	ErrInvalidArabicNum = errors.New("invalid arabic number")
)

var romanToInt = make(map[string]int)
var intToRoman = make(map[int]string)

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

	intToRoman = map[int]string{
		1:  "I",
		2:  "II",
		3:  "III",
		4:  "IV",
		5:  "V",
		6:  "VI",
		7:  "VII",
		8:  "VIII",
		9:  "IX",
		10: "X",
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

func IntToRoman(num int) (string, error) {
	romanStrNum, ok := intToRoman[num]
	if !ok {
		return "", ErrInvalidArabicNum
	}

	return romanStrNum, nil
}
