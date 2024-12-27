package client

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	order "order_service/generated/order"
	wallet "order_service/generated/wallet"
	"order_service/internal/config"
	"order_service/internal/types"
	"sync"
)

var (
	dbManagerClientInstance *DBManagerClient
	dbManagerOnce           sync.Once
)

type DBManagerClient struct {
	conn    *grpc.ClientConn
	Wallets wallet.WalletsClient
	orders  order.OrderClient
}

func NewDBManagerClient() (*DBManagerClient, error) {
	conn, err := grpc.NewClient(
		config.GetConfig().DBManagerAddress, grpc.WithTransportCredentials(
			insecure.
				NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DBManager service")
		return nil, err
	}

	walletsClient := wallet.NewWalletsClient(conn)
	ordersClient := order.NewOrderClient(conn)

	return &DBManagerClient{
		conn:    conn,
		Wallets: walletsClient,
		orders:  ordersClient,
	}, nil
}

func GetDBManagerClient() (*DBManagerClient, error) {
	var err error
	dbManagerOnce.Do(
		func() {
			dbManagerClientInstance, err = NewDBManagerClient()
		},
	)
	return dbManagerClientInstance, err
}

func (c *DBManagerClient) UpdateOrder(updateOrder *types.Order) error {
	var protoOrder *order.OrderDetails
	protoOrder = updateOrder.ToProto()

	_, err := c.orders.UpdateOrder(context.Background(), protoOrder)
	if err != nil {
		return err
	}

	return nil
}

func (c *DBManagerClient) UpdateWallet(updateWallet *types.Wallet) error {
	var protoWallet *wallet.Wallet
	protoWallet = updateWallet.ToProto()

	_, err := c.Wallets.UpdateWallet(context.Background(), protoWallet)
	if err != nil {
		return err
	}

	return nil
}

func (c *DBManagerClient) Close() error {
	return c.conn.Close()
}
