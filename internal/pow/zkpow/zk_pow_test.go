package zkpow

import (
	"strings"
	"testing"
	"time"
)

func TestGenChallenge(t *testing.T) {
	tests := []struct {
		name       string
		difficulty int
		secret     string
		wantLen    int
		wantDur    time.Duration
	}{
		{
			name:       "simple case",
			difficulty: 2,
			secret:     "mysecret",
			wantLen:    64, // Length of SHA-256 hash in hexadecimal
			wantDur:    5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zk := NewZKPoW(tt.difficulty, tt.secret)
			gotChallenge, diff, gotDur := zk.GenChallenge()
			if diff != tt.difficulty {
				t.Errorf("GenChallenge() got difficult = %d, want %d", diff, tt.difficulty)
			}

			if len(gotChallenge) != tt.wantLen {
				t.Errorf("GenChallenge() got length = %d, want %d", len(gotChallenge), tt.wantLen)
			}
			if gotDur != tt.wantDur {
				t.Errorf("GenChallenge() got duration = %v, want %v", gotDur, tt.wantDur)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		difficulty int
		secret     string
		challenge  string
		nonce      string
		proof      string
		want       bool
	}{
		{
			name:       "valid proof",
			difficulty: 1,
			secret:     "mysecret",
			challenge:  "3d4c4a24b8ea2438f9b39a6d5f1c6f7a8f8f7b3a5b3c4a6f6e4f7a5c8e4f6f7a",
			nonce:      "1",
			proof:      "", // will be computed in the test
			want:       true,
		},
		{
			name:       "invalid proof",
			difficulty: 1,
			secret:     "mysecret",
			challenge:  "3d4c4a24b8ea2438f9b39a6d5f1c6f7a8f8f7b3a5b3c4a6f6e4f7a5c8e4f6f7a",
			nonce:      "1",
			proof:      "fffe1f6f7a8f8f7b3a5b3c4a6f6e4f7a5c8e4f6f7a1",
			want:       false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zk := NewZKPoW(tt.difficulty, tt.secret)

			if i == 0 { // Generate proof for the valid case
				nonce, proof := FindProof(tt.challenge, tt.difficulty)
				tt.nonce = nonce
				tt.proof = proof
			}

			got := zk.Validate(tt.challenge, tt.proof, tt.nonce)

			if got != tt.want {
				t.Errorf("Validate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindProof(t *testing.T) {
	tests := []struct {
		name       string
		difficulty int
		secret     string
		challenge  string
	}{
		{
			name:       "find proof with difficulty 1",
			difficulty: 1,
			secret:     "mysecret",
			challenge:  "3d4c4a24b8ea2438f9b39a6d5f1c6f7a8f8f7b3a5b3c4a6f6e4f7a5c8e4f6f7a",
		},
		{
			name:       "find proof with difficulty 2",
			difficulty: 2,
			secret:     "anothersecret",
			challenge:  "f4a4c3d2b1ea2438f9b39a6d5f1c6f7a8f8f7b3a5b3c4a6f6e4f7a5c8e4f6f7a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, proof := FindProof(tt.challenge, tt.difficulty)
			if !strings.HasPrefix(proof, strings.Repeat("0", tt.difficulty)) {
				t.Errorf("FindProof() proof = %s, does not meet difficulty %d", proof, tt.difficulty)
			}
			if len(nonce) == 0 || len(proof) == 0 {
				t.Errorf("FindProof() generated empty nonce or proof")
			}
		})
	}
}
