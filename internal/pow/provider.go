package pow

import "time"

type ProviderPOW interface {
	GenChallenge() (string, int, time.Duration)
	Validate(challenge, proof, nonce string) bool
}
