package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
	"time"
)

var pubSubLogger = util.GetLogger("pub_sub")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		pubSubLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	subscribeClient1, err := c.Subscribe(context.Background(),
		&pb.SubscribeRequest{Topic: []byte("hhh")})
	go func() {
		for {
			message, err := subscribeClient1.Recv()
			pubSubLogger.Printf("subscribeClient1 message: %+v, err: %+v", message, err)
		}
	}()
	subscribeClient2, err := c.Subscribe(context.Background(),
		&pb.SubscribeRequest{Topic: []byte("hhhaaa")})
	go func() {
		for {
			message, err := subscribeClient2.Recv()
			pubSubLogger.Printf("subscribeClient2 message: %+v, err: %+v", message, err)
		}
	}()

	for i := 0; i < 2; i++ {
		publishResponse, err := c.Publish(context.Background(),
			&pb.PublishRequest{Topic: []byte("hhh"), Payload: []byte("payload" + strconv.Itoa(i))})
		pubSubLogger.Printf("publishResponse: %+v, err: %+v", publishResponse, err)
		publishResponse, err = c.Publish(context.Background(),
			&pb.PublishRequest{Topic: []byte("hhhaaa"),
				Payload: []byte("payloadaaa" + strconv.Itoa(i))})
		pubSubLogger.Printf("publishResponse: %+v, err: %+v", publishResponse, err)
		time.Sleep(1 * time.Second)
	}
}
