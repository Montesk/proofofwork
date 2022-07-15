package pow

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/Montesk/proofofwork/errors"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	ErrHashEncode = errors.String("hash encode error")
)

const (
	nonceMin = 1
	nonceMax = 10

	hashDelimiter = "|"
)

type (
	POW interface {
		Generate(clientId string) (string, error)
		Prove(clientId, hash string) (success bool)
	}

	block struct {
		hash        string
		generatedAt int64
	}

	pow struct {
		mu      *sync.Mutex
		storage map[string]block
	}
)

func New() POW {
	return &pow{
		mu:      new(sync.Mutex),
		storage: map[string]block{},
	}
}

// Generate generates and save in sha-1 [clientId|tms|nonce]
// returns challenge [clientId|tms]
func (p *pow) Generate(clientId string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nonce := randomFrom(nonceMin, nonceMax)
	generatedAt := time.Now().Unix()

	// [id|tms|nonce]
	pattern := strings.Split(fmt.Sprintf("%s%s%d%s%d", clientId, hashDelimiter, generatedAt, hashDelimiter, nonce), hashDelimiter)

	hash := prepareHash(pattern)

	bl := block{
		hash:        hash,
		generatedAt: generatedAt,
	}

	challenge := strings.Split(hash, hashDelimiter)
	if len(challenge) < 2 {
		return "", ErrHashEncode
	}

	p.storage[clientId] = bl

	// [id|tms]
	return strings.Join(challenge[:2], hashDelimiter), nil
}

func (p *pow) Prove(clientId, hash string) (success bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	bl, found := p.storage[clientId]
	if !found {
		return false
	}

	return bl.hash == hash
}

func randomFrom(from, to int) int {
	return rand.Intn(to-from) + from
}

func prepareHash(pattern []string) string {
	result := ""
	for idx, part := range pattern {
		sha := sha1.New()

		sha.Write([]byte(part))

		result += hex.EncodeToString(sha.Sum(nil))

		if idx+1 != len(pattern) {
			result += hashDelimiter
		}
	}

	return result
}

func intSliceContains(needle int, search []int) bool {
	for _, v := range search {
		if v == needle {
			return true
		}
	}

	return false
}
