package middleware

import (
	"context"

	"github.com/PicoTools/pico/pkg/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientInterceptor returns unary interceptor with operator's access header
func UnaryClientInterceptor(t string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, shared.GrpcAuthOperatorHeader, t)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns stream interceptor with operator's access header
func StreamClientInterceptor(t string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, shared.GrpcAuthOperatorHeader, t)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
