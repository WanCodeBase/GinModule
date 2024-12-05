package util

import (
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0->max-min
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		idx := rand.Intn(k)
		c := alphabet[idx]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return randomString(6)
}

func RandomMoney() int64 {
	return randomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP"}
	k := len(currencies)
	return currencies[rand.Intn(k)]
}