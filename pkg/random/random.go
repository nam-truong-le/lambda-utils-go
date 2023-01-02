package random

import (
	"math/rand"
	"time"

	"github.com/samber/lo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func String(size int, charset []rune) string {
	return lo.RandomString(size, charset)
}
