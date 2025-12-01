package shortener

import (
	"math/rand"
	"time"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortCodeLength = 6

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateShortCode() string {
	b := make([]byte, shortCodeLength)

	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}
