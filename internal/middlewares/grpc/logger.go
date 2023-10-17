package grpc

import (
	"context"

	"github.com/Orendev/shortener/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Logger(opts []grpc.ServerOption) []grpc.ServerOption {
	opts = append(
		opts,
		grpc.ChainUnaryInterceptor(func(ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler) (resp interface{}, err error) {

			logger.Log.Info("got incoming GRPC request and response",
				zap.Any("reg", req),
				zap.Any("Full Method", info.FullMethod),
			)

			return handler(ctx, req)
		}),
	)

	return opts
}
