package assets

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

var (
	//go:embed data/quotes.txt
	quotes string
)

type quote struct {
	quotesList []string
}

type Quote interface {
	RandomQuite() string
	GetQuotes() []string
}

func NewQuote() Quote {
	return &quote{quotesList: strings.Split(quotes, "\n")}
}

func (s *quote) RandomQuite() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return s.quotesList[r.Int63n(int64(len(s.quotesList)))]
}

func (s *quote) GetQuotes() []string {
	return s.quotesList
}
