package pow

import (
	"crypto/sha1"
	"encoding/hex"
)

const (
	NonceMin = 1
	NonceMax = 10

	HashDelimiter = "|"
)

type POW interface {
	Generate(clientId string) (string, error)
	Prove(clientId, hash string) (success bool)
}

func Encode(pattern []string) string {
	result := ""
	for idx, part := range pattern {
		sha := sha1.New()

		sha.Write([]byte(part))

		result += hex.EncodeToString(sha.Sum(nil))

		if idx+1 != len(pattern) {
			result += HashDelimiter
		}
	}

	return result
}

func SliceContains[T comparable](needle T, search []T) bool {
	for _, v := range search {
		if v == needle {
			return true
		}
	}

	return false
}
