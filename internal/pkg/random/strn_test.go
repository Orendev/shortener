package random

import "testing"

func TestStrn(t *testing.T) {
	n := 8
	if str := Strn(n); len(str) != n {
		t.Errorf("string length is expected %s is equal to %d", str, n)
	}
}
