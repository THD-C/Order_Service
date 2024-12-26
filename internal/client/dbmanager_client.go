package client

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	order "order_service/generated/order"
	wallet "order_service/generated/wallet"
)

type DBManagerClient struct {
	conn    *grpc.ClientConn
	wallets wallet.WalletsClient
	orders  order.OrderClient
}

func NewDBManagerClient(address string) (*DBManagerClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DBManager service")
		return nil, err
	}

	walletsClient := wallet.NewWalletsClient(conn)
	ordersClient := order.NewOrderClient(conn)

	return &DBManagerClient{
		conn:    conn,
		wallets: walletsClient,
		orders:  ordersClient,
	}, nil
}

func (c *DBManagerClient) Close() error {
	return c.conn.Close()
}
