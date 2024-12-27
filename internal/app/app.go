package app

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"order_service/generated/order"
	"order_service/generated/wallet"
	"order_service/internal/cache"
	"order_service/internal/config"
	"order_service/internal/interceptor"
	"order_service/internal/logger"
	"order_service/internal/server"
)

func startPrometheusMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":8080", nil)
}

func startGRPCServer() {
	log := logger.GetLogger()

	conf := config.GetConfig()
	addr := fmt.Sprintf("%s:%s", conf.Addr, conf.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	grpcMetrics := grpc_prometheus.NewServerMetrics()
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpcMetrics.UnaryServerInterceptor(),
			interceptor.UnaryInterceptor(logger.GetLogger()),
		),
	}

	s := grpc.NewServer(opts...)
	wallet.RegisterWalletsServer(s, &server.WalletServer{})
	order.RegisterOrderServer(s, server.NewOrderServer())
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

func App() {
	logger.Init()
	log := logger.GetLogger()

	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msg("Failed to read config from file")
	}

	err = cache.FetchAllWalletsFromService()
	if err != nil {
		log.Fatal().Msg("Failed to fetch wallets")
	}
	go startPrometheusMetrics()
	startGRPCServer()
}
