package main

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		log.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// 连接服务器
	c := pb.NewSDBClient(conn)
	subscribeClient1, err := c.Subscribe(context.Background(),
		&pb.SubscribeRequest{Topic: []byte("hhh")})
	go func() {
		for {
			message, err := subscribeClient1.Recv()
			log.Printf("subscribeClient1 message: %+v, err: %+v", message, err)
		}
	}()
	subscribeClient2, err := c.Subscribe(context.Background(),
		&pb.SubscribeRequest{Topic: []byte("hhhaaa")})
	go func() {
		for {
			message, err := subscribeClient2.Recv()
			log.Printf("subscribeClient2 message: %+v, err: %+v", message, err)
		}
	}()

	for i := 0; i < 2; i++ {
		publishResponse, err := c.Publish(context.Background(),
			&pb.PublishRequest{Topic: []byte("hhh"), Payload: []byte("payload" + strconv.Itoa(i))})
		log.Printf("publishResponse: %+v, err: %+v", publishResponse, err)
		publishResponse, err = c.Publish(context.Background(),
			&pb.PublishRequest{Topic: []byte("hhhaaa"),
				Payload: []byte("payloadaaa" + strconv.Itoa(i))})
		log.Printf("publishResponse: %+v, err: %+v", publishResponse, err)
		time.Sleep(1 * time.Second)
	}
}
