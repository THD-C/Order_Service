package server

import (
	"context"
	proto "order_service/generated/wallet"
	"order_service/internal/service"
)

type WalletServer struct {
	handler service.WalletService
	proto.UnimplementedWalletsServer
}

func (s *WalletServer) CreateWallet(ctx context.Context, req *proto.Wallet) (
	*proto.Wallet,
	error,
) {
	return s.handler.CreateWallet(ctx, req)
}

func (s *WalletServer) UpdateWallet(ctx context.Context, req *proto.Wallet) (
	*proto.Wallet,
	error,
) {
	return s.handler.UpdateWallet(ctx, req)
}

func (s *WalletServer) GetWallet(ctx context.Context, req *proto.WalletID) (
	*proto.Wallet,
	error,
) {
	return s.handler.GetWallet(ctx, req)
}
