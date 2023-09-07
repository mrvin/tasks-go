package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mrvin/tasks-go/go-calc/calc"
)

const (
	arabicMode uint8 = 0
	romanMode  uint8 = 2
)

func main() {
	var expressionStr string
	fmt.Print("Input: ")
	expressionStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Printf("go-calc: can't read string: %v", err)
		return
	}
	expressionSl := strings.Split(expressionStr[:len(expressionStr)-1], " ")

	if len(expressionSl) != 3 {
		log.Print("go-calc: should be 3 words")
		return
	}

	var mode uint8
	num1, err := strconv.Atoi(expressionSl[0])
	if err != nil {
		num1, err = calc.RomanToInt(expressionSl[0])
		if err != nil {
			log.Printf("go-calc: can't convert: %v", err)
			return
		}
		mode++
	}
	num2, err := strconv.Atoi(expressionSl[2])
	if err != nil {
		num2, err = calc.RomanToInt(expressionSl[2])
		if err != nil {
			log.Printf("go-calc: can't convert: %v", err)
			return
		}
		mode++
	}

	// validation
	if num1 < 1 || num1 > 10 || num2 < 1 || num2 > 10 {
		log.Print("go-calc: numbers are out of range [1;10]")
		return
	}
	if mode != arabicMode && mode != romanMode {
		log.Print("go-calc: different number systems")
		return
	}

	result, err := calc.Calc(num1, num2, expressionSl[1])
	if err != nil {
		log.Printf("go-calc: %v", err)
		return
	}

	switch mode {
	case arabicMode:
		fmt.Printf("result: %d\n", result)
	case romanMode:
		resultRoman, err := calc.IntToRoman(result)
		if err != nil {
			log.Printf("go-calc: result %v: %d", err, result)
			return
		}
		fmt.Printf("result: %s\n", resultRoman)
	}
}
