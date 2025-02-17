package service

import (
	operatorv1 "github.com/PicoTools/pico/pkg/proto/operator/v1"
	"github.com/lrita/cmap"
	"google.golang.org/grpc"
)

type streams struct {
	controlStream  grpc.ServerStreamingClient[operatorv1.HelloResponse]
	tasksStream    grpc.BidiStreamingClient[operatorv1.SubscribeTasksRequest, operatorv1.SubscribeTasksResponse]
	commandStreams cmap.Map[uint32, grpc.ClientStreamingClient[operatorv1.NewCommandRequest, operatorv1.NewCommandResponse]]
}
