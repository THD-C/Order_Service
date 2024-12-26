package order_service

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"order_service/internal/cache"
	"order_service/internal/types"
)

func FetchAndValidateWallets(req *types.Order) (*types.Wallet, *types.Wallet, error) {
	fiatWallet, err := cache.FetchWallet(req.FiatWalletID)
	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to get fiat wallet")
		return nil, nil, fmt.Errorf("failed to get fiat wallet: %v", err)
	}

	cryptoWallet, err := cache.FetchWallet(req.CryptoWalletID)
	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to get crypto wallet")
		return nil, nil, fmt.Errorf("failed to get crypto wallet: %v", err)
	}

	return fiatWallet, cryptoWallet, nil
}

func RollbackWallet(
	wallet *types.Wallet,
	originalValue decimal.Decimal,
) error {
	wallet.Value = originalValue
	if err := cache.SaveWallet(wallet); err != nil {
		log.Error().Err(err).Msg("Rollback failed")
		return err
	}
	return nil
}
