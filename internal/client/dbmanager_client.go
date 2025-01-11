package client

import (
	"context"
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

func (c *DBManagerClient) FetchAllPendingOrders() ([]*types.Order, error) {
	orders := make([]*types.Order, 0)

	ordersList, err := c.orders.GetOrders(
		context.Background(),
		&order.OrderFilter{
			Type:   order.OrderType_ORDER_TYPE_PENDING,
			Status: order.OrderStatus_ORDER_STATUS_PENDING,
		},
	)
	if err != nil {
		return nil, err
	}

	for _, protoOrder := range ordersList.Orders {
		var myOrder types.Order
		if err := myOrder.FromProto(protoOrder); err != nil {
			continue
		}
		orders = append(orders, &myOrder)
	}

	return orders, nil
}

func (c *DBManagerClient) FetchAllWallets() ([]*types.Wallet, error) {
	wallets := make([]*types.Wallet, 0)

	walletList, err := c.Wallets.GetAllWallets(context.Background(), &wallet.Empty{})
	if err != nil {
		return nil, err
	}

	for _, protoWallet := range walletList.Wallets {
		var myWallet types.Wallet
		if err := myWallet.FromProto(protoWallet); err != nil {
			continue
		}
		wallets = append(wallets, &myWallet)
	}

	return wallets, nil
}

func (c *DBManagerClient) Close() error {
	return c.conn.Close()
}
