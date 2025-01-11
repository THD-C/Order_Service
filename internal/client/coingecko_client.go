package client

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "order_service/generated/coins"
	"order_service/internal/config"
	"order_service/internal/logger"
	"order_service/internal/types"
	"strings"
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

func (c *CoinGeckoClient) GetAllCoinsPrice() ([]*types.CoinPrice, error) {
	conf := config.GetConfig()
	log := logger.GetLogger()
	ctx, cancel := context.WithTimeout(context.Background(), conf.CoingeckoServiceTimeout)
	defer cancel() // Ensure the context is canceled to avoid resource leaks

	res, err := c.client.GetAllCoinsPrices(ctx, &proto.AllCoinsPricesRequest{})
	if err != nil {
		return nil, err
	}

	coinPrices := make([]*types.CoinPrice, 0)
	for coinID, coinData := range res.Data.AsMap() {
		coinMap, ok := coinData.(map[string]interface{})
		if !ok {
			log.Error().Msgf("invalid coin data for %s", coinID)
			continue
		}

		prices := make(map[string]decimal.Decimal)
		for currency, value := range coinMap {
			decimalValue, err := decimal.NewFromString(fmt.Sprintf("%v", value))
			if err != nil {
				log.Error().Msgf("invalid price for %s in %s: %v", coinID, currency, err)
				continue
			}
			prices[strings.ToUpper(currency)] = decimalValue
		}

		coinPrices = append(
			coinPrices, &types.CoinPrice{
				CoinId: coinID,
				Prices: prices,
			},
		)
	}

	return coinPrices, nil
}

func (c *CoinGeckoClient) Close() error {
	return c.conn.Close()
}
