package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	conf := config.GetConfig()
	addr := fmt.Sprintf(":%s", conf.PrometheusPort)
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(addr, nil)
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
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	s := grpc.NewServer(opts...)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	wallet.RegisterWalletsServer(s, &server.WalletServer{})
	order.RegisterOrderServer(s, server.NewOrderServer())
	reflection.Register(s)

	log.Info().Msgf("server listening at %v", lis.Addr())
	healthServer.SetServingStatus(conf.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	if err = s.Serve(lis); err != nil {
		healthServer.SetServingStatus(conf.ServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}

func populateCaches() error {
	var err error
	conf := config.GetConfig()

	priceCache := cache.NewPriceCache()
	go priceCache.UpdatePrices(context.Background(), conf.PendingOrderCheckFrequency)

	orderCache := cache.GetOrderCache()
	err = orderCache.FetchAllOrders()
	if err != nil {
		return fmt.Errorf("failed to fetch orders from database")
	}

	err = cache.FetchAllWalletsFromService()
	if err != nil {
		return fmt.Errorf("failed to fetch wallets from database")
	}

	return nil
}

func Init() error {
	logger.Init()
	log := logger.GetLogger()

	_, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msg("Failed to read config from file")
		return err
	}

	err = populateCaches()
	if err != nil {
		log.Fatal().Err(err)
		return err
	}

	return nil
}

func Run() {
	log := logger.GetLogger()

	go startPrometheusMetrics()

	tp := config.Init()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	startGRPCServer()
}
