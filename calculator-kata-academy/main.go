package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mrvin/tasks-go/calculator-kata-academy/calc"
)

func main() {
	if len(os.Args) > 4 {
		log.Print("go-calc: too many arguments")
		return
	}
	if len(os.Args[1:]) != 3 {
		log.Print("go-calc: not enough arguments < 3")
		return
	}
	num1, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Printf("go-calc: can't convert: %v", err)
		num1, err = calc.RomanToInt(os.Args[1])
		if err != nil {
			log.Printf("go-calc: can't convert: %v", err)
			return
		}
	}

	num2, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Printf("go-calc: can't convert: %v", err)
		num2, err = calc.RomanToInt(os.Args[3])
		if err != nil {
			log.Printf("go-calc: can't convert: %v", err)
			return
		}
	}

	// validation
	if num1 < 1 || num1 > 10 || num2 < 1 || num2 > 10 {
		log.Print("go-calc: numbers are out of range [1;10]")
		return
	}

	result, err := calc.Calc(num1, num2, os.Args[2])
	if err != nil {
		log.Printf("go-calc: %v", err)
		return
	}
	fmt.Println(result)
}
