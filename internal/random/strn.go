package random

import (
	"math/rand"
	"strings"
	"time"
)

func Strn(n int) string {
	rand.Seed(time.Now().UnixNano())
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	chars := []rune(alphabet + strings.ToLower(alphabet))

	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}
