package cache

import (
	"fmt"
	"github.com/shopspring/decimal"
	"order_service/internal/types"
	"sync"
)

var (
	walletMap sync.Map
)

func SaveWallet(wallet *types.Wallet) error {
	walletMap.Store(wallet.ID, wallet)
	return nil
}

func FetchWallet(walletID string) (*types.Wallet, error) {
	value, ok := walletMap.Load(walletID)
	if !ok {
		return nil, fmt.Errorf("wallet not found")
	}

	wallet, ok := value.(*types.Wallet)
	if !ok {
		return nil, fmt.Errorf("failed to assert wallet data")
	}

	return wallet, nil
}

func FetchAllWalletsFromService() error {
	// TODO: Docelowo z serwisu Stacha to będzie pobierane,
	//  ale na razie nie ma tam danych więc sam symuluje

	wallets, err := fetchWalletsFromExternalService()
	if err != nil {
		return fmt.Errorf("failed to fetch wallets from service: %v", err)
	}

	for _, wallet := range wallets {
		walletMap.Store(wallet.ID, wallet)
	}

	return nil
}

func fetchWalletsFromExternalService() ([]*types.Wallet, error) {
	return []*types.Wallet{
		{ID: "1", Currency: "USD", Value: decimal.NewFromFloat(100), UserID: "1", IsCrypto: false},
		{ID: "2", Currency: "BTC", Value: decimal.NewFromFloat(0), UserID: "1", IsCrypto: true},
	}, nil
}
