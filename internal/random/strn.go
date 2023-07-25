package random

import (
	"math/rand"
	"strings"
	"time"
)

// Strn specify the length of the random string consisting of letters.
func Strn(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	chars := []rune(alphabet + strings.ToLower(alphabet))

	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteRune(chars[r.Intn(len(chars))])
	}

	return b.String()
}
