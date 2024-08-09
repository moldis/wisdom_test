package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/moldis/wisdom_test/assets"
	"github.com/moldis/wisdom_test/config"
	"github.com/moldis/wisdom_test/internal/pow"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"time"
)

type TcpServer struct {
	cfg    *config.Config
	pow    pow.ProviderPOW
	quotes assets.Quote
}

func NewTcpServer(cfg *config.Config, pow pow.ProviderPOW, quotes assets.Quote) *TcpServer {
	return &TcpServer{cfg: cfg, pow: pow, quotes: quotes}
}

func (s *TcpServer) Start(ctx context.Context) {
	lAddr, err := net.ResolveTCPAddr("tcp", s.cfg.Addr)
	if err != nil {
		log.Fatal().Msgf("resolve local address: %s", err)
	}

	l, err := net.ListenTCP("tcp", lAddr)
	if err != nil {
		log.Fatal().Msgf("can't start listener: %s", err)
	}

	defer func() {
		err = l.Close()
		if err != nil {
			log.Err(err).Msgf("can't close listener conn: %s", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("context is done, exiting")
			return
		default:
			c, acceptErr := l.Accept()
			if acceptErr != nil {
				log.Err(acceptErr).Msgf("client can't Accept")
				continue
			}

			go s.handleClient(ctx, c)
		}
	}
}

func (s *TcpServer) handleClient(ctx context.Context, c net.Conn) {
	defer func(c net.Conn) {
		if err := c.Close(); err != nil {
			log.Err(err).Msgf("can't close client conn")
		}
	}(c)

	challenge, diff, deadline := s.pow.GenChallenge()
	_, err := io.WriteString(c, fmt.Sprintf("%s %d\n", challenge, diff))
	if err != nil {
		log.Err(err).Msgf("can't write challenge to client, closing")
		return
	}

	if err = c.SetReadDeadline(time.Now().Add(deadline)); err != nil {
		log.Err(err).Msgf("deadline exceed, client can't resolve challenge")
		return
	}

	reader := bufio.NewReader(c)
	proof, err := reader.ReadString('\n')
	if err != nil {
		log.Err(err).Msgf("Error reading proof")
		return
	}
	proof = proof[:len(proof)-1]

	nonce, err := reader.ReadString('\n')
	if err != nil {
		log.Err(err).Msgf("error reading nonce")
		return
	}
	nonce = nonce[:len(nonce)-1]

	valid := s.pow.Validate(challenge, proof, nonce)
	if !valid {
		log.Error().Msgf("client invalid nonce from client")
		return
	}

	quote := s.quotes.RandomQuite()
	_, err = c.Write([]byte(quote + "\n"))
	if err != nil {
		log.Error().Msgf("can't write quote to client, closing")
		return
	}

	log.Info().Msgf("quote succesfully sent: %s", quote)
}
