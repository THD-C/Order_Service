package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "order_service/generated/coins"
	"order_service/internal/config"
	"sync"
)

var (
	coinGeckoClientInstance *CoinGeckoClient
	coinGeckoOnce           sync.Once
)

type CoinGeckoClient struct {
	conn   *grpc.ClientConn
	client proto.CoinsClient
}

func NewCoinGeckoClient() (*CoinGeckoClient, error) {
	conn, err := grpc.NewClient(
		config.GetConfig().CoingeckoServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := proto.NewCoinsClient(conn)
	return &CoinGeckoClient{
		conn:   conn,
		client: client,
	}, nil
}

func GetCoinGeckoClient() (*CoinGeckoClient, error) {
	var err error
	coinGeckoOnce.Do(
		func() {
			coinGeckoClientInstance, err = NewCoinGeckoClient()
		},
	)
	return coinGeckoClientInstance, err
}

func (c *CoinGeckoClient) Close() error {
	return c.conn.Close()
}
