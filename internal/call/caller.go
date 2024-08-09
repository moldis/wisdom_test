package call

import (
	"bufio"
	"context"
	"fmt"
	"github.com/moldis/wisdom_test/internal/pow/zkpow"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

type Caller struct {
	protocol string
	addr     string
}

func NewCaller(optFns ...OptFn) *Caller {
	c := &Caller{}

	for _, fn := range optFns {
		fn(c)
	}

	return c
}

func (c *Caller) Run(ctx context.Context) error {
	conn, err := net.Dial(c.protocol, c.addr)
	if err != nil {
		return fmt.Errorf("failed to dial the server: %w", err)
	}

	var challenge string
	var difficulty int
	if _, err = fmt.Fscanf(conn, "%s %d\n", &challenge, &difficulty); err != nil {
		return fmt.Errorf("can't read challenge and diff")
	}

	nonce, proof := zkpow.FindProof(challenge, difficulty)

	message := fmt.Sprintf("%s\n%s\n", proof, nonce)
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return fmt.Errorf("failed to set a reading deadline: %w", err)
	}

	reader := bufio.NewReader(conn)
	quote, err := reader.ReadString('\n')
	if err != nil {
		log.Err(err).Msgf("Error reading quote")
		return err
	}

	log.Info().Msgf("received quote: %s", quote)

	return nil
}
