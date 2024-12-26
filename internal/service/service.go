package service

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/PicoTools/pico-cli/internal/middleware"
	operatorv1 "github.com/PicoTools/pico-shared/proto/gen/operator/v1"
	"github.com/PicoTools/pico-shared/shared"
	"github.com/go-faster/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var operatorConn = &grpcConn{}

type grpcConn struct {
	ctx      context.Context
	conn     *grpc.ClientConn
	ss       streams
	metadata metadata
	svc      operatorv1.OperatorServiceClient
}

// Init initializes operator's GRPC client
func Init(ctx context.Context, host string, token string) error {
	var err error
	operatorConn.ctx = ctx

	if operatorConn.conn, err = grpc.NewClient(
		host,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
		grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor(token)),
		grpc.WithStreamInterceptor(middleware.StreamClientInterceptor(token)),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(shared.MaxProtobufMessageSize),
			grpc.MaxCallSendMsgSize(shared.MaxProtobufMessageSize),
		),
	); err != nil {
		return err
	}
	operatorConn.svc = operatorv1.NewOperatorServiceClient(operatorConn.conn)

	// open connection, authenticate and get server data
	operatorConn.ss.controlStream, err = HelloInit(ctx)
	if err != nil {
		return errors.Wrap(err, "open hello stream")
	}
	if err = HelloHandshake(ctx); err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return fmt.Errorf("%s", st.Message())
			default:
				return errors.Wrap(err, "process hello handshake")
			}
		}
		return errors.Wrap(err, "process hello handshake")
	}

	// start subscriptions
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return HelloMonitor(ctx)
	})
	g.Go(func() error {
		return SubscribeChat(ctx)
	})
	g.Go(func() error {
		return SubscribeAnts(ctx)
	})
	g.Go(func() error {
		return SubscribeTasks(ctx)
	})

	go func() {
		g.Wait()
	}()

	return nil
}

// Close closes operator's GRPC connection
func Close() error {
	if operatorConn.conn != nil {
		return operatorConn.conn.Close()
	}
	return nil
}

func getSvc() operatorv1.OperatorServiceClient {
	return operatorConn.svc
}
