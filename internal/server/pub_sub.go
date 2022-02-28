package server

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
	"time"
)

type PubSubServer struct {
	pb.UnimplementedSDBServer
}

func (server *PubSubServer) Subscribe(request *pb.SubscribeRequest, subscribeServer pb.SDB_SubscribeServer) error {
	_, err := service.Subscribe(request.Topic, &subscribeServer)
	if err != nil {
		return err
	}
	for {
		ch := service.GetStopChannel(&subscribeServer)
		select {
		// stop
		case <-ch:
			service.CleanSubscribeServer(&subscribeServer)
			return nil
		case <-time.After(5 * time.Second):
			continue
		}
	}
}

func (server *PubSubServer) Publish(_ context.Context, request *pb.PublishRequest) (*pb.PublishResponse, error) {
	res, err := service.Publish(request)
	return &pb.PublishResponse{Success: res}, err
}
