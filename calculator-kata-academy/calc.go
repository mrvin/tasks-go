package main

func Calc(a, b int, operation rune) (int, error) {
	switch operation {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		return a / b, nil
	}

	return 0, ErrInvalidOperation
}
