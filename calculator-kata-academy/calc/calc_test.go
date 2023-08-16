package calc

import (
	"testing"
)

func TestCalc(t *testing.T) {
	var tests = []struct {
		a, b      int
		operation string
		want      int
		err       error
	}{
		{1, 2, "+", 3, nil},
		{3, 3, "+", 6, nil},
		{5, 4, "+", 9, nil},
		{3, 4, "%", 0, ErrInvalidOperation},
	}

	for _, test := range tests {
		if got, err := Calc(test.a, test.b, test.operation); err != test.err || got != test.want {
			t.Errorf("Calc(%d, %d, %q) = %d", test.a, test.b, test.operation, got)
		}
	}
}

func TestRomanToInt(t *testing.T) {
	var tests = []struct {
		romanStrNum string
		want        int
		err         error
	}{
		{"X", 10, nil},
		{"IX", 9, nil},
		{"VIII", 8, nil},
		{"VII", 7, nil},
		{"VI", 6, nil},
		{"V", 5, nil},
		{"IV", 4, nil},
		{"III", 3, nil},
		{"II", 2, nil},
		{"I", 1, nil},
		{"XXX", 0, ErrInvalidRomanNum},
	}

	for _, test := range tests {
		if got, err := RomanToInt(test.romanStrNum); err != test.err || got != test.want {
			t.Errorf("RomanToInt(%q) = %d", test.romanStrNum, got)
		}
	}
}
