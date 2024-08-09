package assets

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomQuote(t *testing.T) {
	q := NewQuote()
	randomQuite := q.RandomQuite()
	require.NotEmpty(t, randomQuite)

	found := false
	for _, q := range q.GetQuotes() {
		if randomQuite == q {
			found = true
			break
		}
	}

	require.True(t, found)
}

func TestRandomQuoteConsistency(t *testing.T) {
	q := NewQuote()
	quote1 := q.RandomQuite()
	quote2 := q.RandomQuite()

	require.NotEqual(t, quote1, quote2)
}
