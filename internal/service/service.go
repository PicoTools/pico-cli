package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/PicoTools/pico-cli/internal/middleware"
	operatorv1 "github.com/PicoTools/pico/pkg/proto/operator/v1"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/fatih/color"
	"github.com/go-faster/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var conn = &grpcConn{}

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
	conn.ctx = ctx

	if conn.conn, err = grpc.NewClient(
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
	conn.svc = operatorv1.NewOperatorServiceClient(conn.conn)

	// open connection, authenticate and get server data
	conn.ss.controlStream, err = HelloInit(ctx)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unavailable:
				return fmt.Errorf("operator's server is unavailable on %s", host)
			}
		}
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

	// monitor GRPC health
	go handleConnectionHealth(ctx, conn.conn)

	// start subscriptions
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return HelloMonitor(ctx)
	})
	g.Go(func() error {
		return SubscribeChat(ctx)
	})
	g.Go(func() error {
		return SubscribeAgents(ctx)
	})
	g.Go(func() error {
		return SubscribeTasks(ctx)
	})

	go func() {
		_ = g.Wait()
	}()

	return nil
}

// handleConnectionHealth check if connection to server lost
func handleConnectionHealth(ctx context.Context, conn *grpc.ClientConn) {
	if conn.WaitForStateChange(ctx, conn.GetState()) {
		newState := conn.GetState()
		if newState == connectivity.Idle {
			fmt.Println(color.RedString("\n\nConnection to operator's server lost. Exiting."))
			os.Exit(-2)
		}
	}
}

// Close closes operator's GRPC connection
func Close() error {
	if conn.conn != nil {
		return conn.conn.Close()
	}
	return nil
}

// getSvc returns service object to interact with GRPC server
func getSvc() operatorv1.OperatorServiceClient {
	return conn.svc
}

// GetMetadata returns metadata of GRPC connection
func (g *grpcConn) getMetadata() metadata {
	return g.metadata
}

// GetUsername returns username of operator
func GetUsername() string {
	return conn.getMetadata().GetUsername()
}
