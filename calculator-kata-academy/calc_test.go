package main

import (
	"testing"
)

func TestCalc(t *testing.T) {
	var tests = []struct {
		a, b      int
		operation rune
		want      int
		err       error
	}{
		{1, 2, '+', 3, nil},
		{3, 3, '+', 6, nil},
		{5, 4, '+', 9, nil},
	}

	for _, test := range tests {
		if got, err := Calc(test.a, test.b, test.operation); err != test.err || got != test.want {
			t.Errorf("Calc(%d, %d, %c) = %d", test.a, test.b, test.operation, got)
		}
	}
}
