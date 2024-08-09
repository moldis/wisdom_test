package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/moldis/wisdom_test/assets"
	"github.com/moldis/wisdom_test/config"
	"github.com/moldis/wisdom_test/internal/pow/zkpow"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func Test_HandleChallenge_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pow := zkpow.NewZKPoW(1, "my_secret")
	quotes := assets.NewQuote()

	server := NewTcpServer(&config.Config{Addr: ":8787", PoWDifficulty: 1}, pow, quotes)
	go server.Start(ctx)

	conn, err := net.Dial("tcp", "localhost:8787")
	require.NoError(t, err)
	defer conn.Close()

	var challenge string
	var difficulty int
	_, err = fmt.Fscanf(conn, "%s %d\n", &challenge, &difficulty)
	require.NoError(t, err)

	nonce, proof := zkpow.FindProof(challenge, difficulty)

	message := fmt.Sprintf("%s\n%s\n", proof, nonce)
	_, err = conn.Write([]byte(message))
	require.NoError(t, err)

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	require.NoError(t, err)

	reader := bufio.NewReader(conn)
	quote, err := reader.ReadString('\n')
	require.NoError(t, err)
	require.NotEmpty(t, quote)

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func Test_HandleChallenge_Failed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pow := zkpow.NewZKPoW(1, "my_secret")
	quotes := assets.NewQuote()

	server := NewTcpServer(&config.Config{Addr: ":8788", PoWDifficulty: 1}, pow, quotes)
	go server.Start(ctx)

	conn, err := net.Dial("tcp", "localhost:8788")
	require.NoError(t, err)
	defer conn.Close()

	var challenge string
	var difficulty int
	_, err = fmt.Fscanf(conn, "%s %d\n", &challenge, &difficulty)
	require.NoError(t, err)

	message := fmt.Sprintf("%s\n%s\n", "proof", "nonce")
	_, err = conn.Write([]byte(message))
	require.NoError(t, err)

	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	require.NoError(t, err)

	reader := bufio.NewReader(conn)
	_, err = reader.ReadString('\n')
	require.Error(t, err)

	cancel()
	time.Sleep(100 * time.Millisecond)
}
