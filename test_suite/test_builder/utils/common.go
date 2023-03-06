package utils

import (
	"math/rand"
	"time"
)

func RandomStringWithCharset(length int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomRange(min, max int) int64 {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(50 * time.Millisecond)
	return int64(rand.Intn(max-min+1) + min)
}
