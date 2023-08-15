package main

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
)

func main() {
	result, err := Calc(2, 2, '+')
	if err != nil {
		log.Printf("Calc: %v", err)
	}
	fmt.Println(result)
}
