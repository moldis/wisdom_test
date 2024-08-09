package zkpow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/moldis/wisdom_test/internal/pow"
	"strings"
	"time"
)

type zk struct {
	difficulty int
	secret     string
}

func NewZKPoW(difficulty int, secret string) pow.ProviderPOW {
	return &zk{difficulty: difficulty, secret: secret}
}

func (s *zk) GenChallenge() (string, int, time.Duration) {
	hash := sha256.Sum256([]byte(s.secret))
	return hex.EncodeToString(hash[:]), s.difficulty, 5 * time.Second
}

func (s *zk) Validate(challenge, proof, nonce string) bool {
	computedProof := ComputeProof(challenge, nonce)
	if computedProof != proof {
		return false
	}
	return strings.HasPrefix(proof, strings.Repeat("0", s.difficulty))
}

func ComputeProof(challenge, nonce string) string {
	proofInput := challenge + nonce
	hash := sha256.Sum256([]byte(proofInput))
	return hex.EncodeToString(hash[:])
}

func FindProof(challenge string, difficulty int) (nonce string, proof string) {
	for i := 0; ; i++ {
		nonce = fmt.Sprintf("%x", i)
		proof = ComputeProof(challenge, nonce)
		if strings.HasPrefix(proof, strings.Repeat("0", difficulty)) {
			return nonce, proof
		}
	}
}
