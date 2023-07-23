package random

import (
	"testing"
)

func TestStrn(t *testing.T) {
	n := 8
	if str := Strn(n); len(str) != n {
		t.Errorf("string length is expected %s is equal to %d", str, n)
	}
}
func BenchmarkStrn(b *testing.B) {
	n := 8
	for i := 0; i < b.N; i++ {
		Strn(n)
	}
}
