package service

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"order_service/generated/wallet"
	"order_service/internal/cache"
	"order_service/internal/logger"
	"order_service/internal/types"
)

type WalletService struct {
}

func (s *WalletService) CreateWallet(_ context.Context, req *wallet.Wallet) (
	*wallet.Wallet,
	error,
) {
	log := logger.GetLogger()
	log.Info().Interface("request", req).Msg("Creating wallet")

	var createWallet types.Wallet
	err := createWallet.FromProto(req)
	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to convert proto to wallet")
		return nil, err
	}

	if _, err = cache.FetchWallet(createWallet.ID); err == nil {
		log.Error().Err(err).Interface("request", req).Msg("Wallet already exists")
		return req, nil
	}

	createWallet.Mutex.Lock()
	err = cache.SaveWallet(&createWallet)
	createWallet.Mutex.Unlock()

	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to save wallet")
		return nil, err
	}

	log.Info().Interface("request", req).Msg("Wallet created successfully")
	return createWallet.ToProto(), nil
}

func (s *WalletService) UpdateWallet(_ context.Context, req *wallet.Wallet) (
	*wallet.Wallet,
	error,
) {
	log := logger.GetLogger()
	log.Info().Interface("request", req).Msg("Updating wallet")

	var updateWallet *types.Wallet
	updateWallet, err := cache.FetchWallet(req.Id)
	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Wallet not found")
		return nil, fmt.Errorf("wallet not found: %v", err)
	}

	updateWallet.Mutex.Lock()
	updateWallet.Value, err = decimal.NewFromString(req.Value)
	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to read new value for wallet")
		return nil, fmt.Errorf("failed to read new value for wallet: %v", err)
	}

	err = cache.SaveWallet(updateWallet)
	updateWallet.Mutex.Unlock()

	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to update wallet")
		return nil, fmt.Errorf("failed to update wallet: %v", err)
	}

	log.Info().Interface("request", req).Msg("Wallet updated successfully")
	return updateWallet.ToProto(), nil
}

func (s *WalletService) GetWallet(_ context.Context, req *wallet.WalletID) (
	*wallet.Wallet,
	error,
) {
	log := logger.GetLogger()
	log.Info().Interface("request", req).Msg("Getting wallet")

	var retrievedWallet *types.Wallet
	retrievedWallet, err := cache.FetchWallet(req.Id)

	if err != nil {
		log.Error().Err(err).Interface("request", req).Msg("Failed to get wallet")
		return nil, fmt.Errorf("failed to get wallet: %v", err)
	}

	log.Info().Interface("request", req).Msg("Wallet retrieved successfully")
	return retrievedWallet.ToProto(), nil
}
