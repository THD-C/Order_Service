package app

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"order_service/generated/order"
	"order_service/generated/wallet"
	"order_service/internal/cache"
	"order_service/internal/config"
	"order_service/internal/server"
)

func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	wallet.RegisterWalletsServer(s, &server.WalletServer{})
	order.RegisterOrderServer(s, server.NewOrderServer())
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

func App() {
	_, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal().Msg("Failed to read config from file")
	}

	err = cache.FetchAllWalletsFromService()
	if err != nil {
		log.Fatal().Msg("Failed to fetch wallets")
	}
	startGRPCServer()
}
