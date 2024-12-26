package client

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	proto "order_service/generated/coins"
)

type CoinGeckoClient struct {
	conn   *grpc.ClientConn
	client proto.CoinsClient
}

func NewCoinGeckoClient(address string) (*CoinGeckoClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to coin gecko service")
		return nil, err
	}

	client := proto.NewCoinsClient(conn)
	return &CoinGeckoClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *CoinGeckoClient) Close() error {
	return c.conn.Close()
}
