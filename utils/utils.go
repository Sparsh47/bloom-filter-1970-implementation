package utils

import (
	"unicode"

	"github.com/spaolacci/murmur3"
)

func Hash(data []byte, seed int) uint32 {
	return murmur3.Sum32WithSeed(data, uint32(seed))
}

func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
