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

func TestIntToRoman(t *testing.T) {
	var tests = []struct {
		num  int
		want string
		err  error
	}{
		{100, "C", nil},
		{90, "XC", nil},
		{80, "LXXX", nil},
		{70, "LXX", nil},
		{60, "LX", nil},
		{50, "L", nil},
		{40, "XL", nil},
		{30, "XXX", nil},
		{28, "XXVIII", nil},
		{27, "XXVII", nil},
		{26, "XXVI", nil},
		{25, "XXV", nil},
		{24, "XXIV", nil},
		{23, "XXIII", nil},
		{22, "XXII", nil},
		{21, "XXI", nil},
		{20, "XX", nil},
		{19, "XIX", nil},
		{18, "XVIII", nil},
		{17, "XVII", nil},
		{16, "XVI", nil},
		{15, "XV", nil},
		{14, "XIV", nil},
		{13, "XIII", nil},
		{12, "XII", nil},
		{11, "XI", nil},
		{10, "X", nil},
		{9, "IX", nil},
		{8, "VIII", nil},
		{7, "VII", nil},
		{6, "VI", nil},
		{5, "V", nil},
		{4, "IV", nil},
		{3, "III", nil},
		{2, "II", nil},
		{1, "I", nil},
		{0, "", ErrInvalidArabicNum},
	}

	for _, test := range tests {
		if got, err := IntToRoman(test.num); err != test.err || got != test.want {
			t.Errorf("RomanToInt(%d) = %q", test.num, got)
		}
	}
}
