package fibonacci

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	var tests = []struct {
		prevNumStr, numStr string
		from, to           uint64
		want               []string
	}{
		{"0", "1", 1, 5, []string{"1", "1", "2", "3", "5"}},
		{"0", "1", 1, 4, []string{"1", "1", "2", "3"}},
		{"1", "2", 3, 5, []string{"2", "3", "5"}},
	}

	for _, test := range tests {
		if got := Get(test.prevNumStr, test.numStr, test.from, test.to); !reflect.DeepEqual(got, test.want) {
			t.Errorf("IsPalindrome(\"%s\", \"%s\", %d, %d) = %v; want = %v", test.prevNumStr, test.numStr, test.from, test.to, got, test.want)
		}
	}
}
