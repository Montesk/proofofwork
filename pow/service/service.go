// Package that provides implementation of mem-cache POW proof and generator mechanism

package service

import (
	"fmt"
	"github.com/Montesk/proofofwork/errors"
	"github.com/Montesk/proofofwork/pow/pow"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	ErrHashEncode = errors.String("hash encode error")
)

type (
	block struct {
		hash        string
		generatedAt int64
	}

	service struct {
		mu      *sync.Mutex
		storage map[string]block
	}
)

func New() pow.POW {
	return &service{
		mu:      new(sync.Mutex),
		storage: map[string]block{},
	}
}

// Generate generates and save in sha-1 [clientId|tms|nonce]
// returns challenge [clientId|tms]
func (p *service) Generate(clientId string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nonce := randomFrom(pow.NonceMin, pow.NonceMax)
	generatedAt := time.Now().Unix()

	// [id|tms|nonce]
	pattern := strings.Split(fmt.Sprintf("%s%s%d%s%d", clientId, pow.HashDelimiter, generatedAt, pow.HashDelimiter, nonce), pow.HashDelimiter)

	hash := pow.Encode(pattern)

	bl := block{
		hash:        hash,
		generatedAt: generatedAt,
	}

	challenge := strings.Split(hash, pow.HashDelimiter)
	if len(challenge) < 2 {
		return "", ErrHashEncode
	}

	p.storage[clientId] = bl

	// [id|tms]
	return strings.Join(challenge[:2], pow.HashDelimiter), nil
}

func (p *service) Prove(clientId, hash string) (success bool) {
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
