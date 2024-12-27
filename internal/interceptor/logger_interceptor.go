package interceptor

import (
	"context"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"time"
)

func UnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		logger.Info().
			Str("method", info.FullMethod).
			Interface("request", req).
			Msg("Incoming gRPC request")

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			logger.Error().
				Str("method", info.FullMethod).
				Interface("request", req).
				Err(err).
				Dur("duration", duration).
				Msg("gRPC request failed")
		} else {
			logger.Info().
				Str("method", info.FullMethod).
				Interface("request", req).
				Interface("response", resp).
				Dur("duration", duration).
				Msg("gRPC request succeeded")
		}

		return resp, err
	}
}
