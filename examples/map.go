package main

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
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
	// 发起 mpush 请求
	pairs := make([]*pb.Pair, 100)
	for i := 0; i < 100; i++ {
		pairs[i] = &pb.Pair{Key: []byte("k" + fmt.Sprint(i)), Value: []byte("v" + fmt.Sprint(i+1))}
	}
	mpushResponse, err := c.MPush(context.Background(),
		&pb.MPushRequest{Key: []byte("h"), Pairs: pairs})
	log.Printf("mpushResponse: %+v, err: %+v", mpushResponse, err)

	mmembersResponse, _ := c.MMembers(context.Background(),
		&pb.MMembersRequest{Key: []byte("h")})
	log.Printf("mmembersResponse: %+v, err: %+v", mmembersResponse, err)

	// 发起 mpop 请求
	keys := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		keys[i] = []byte("k" + fmt.Sprint(i*2))
	}
	mpopResponse, err := c.MPop(context.Background(),
		&pb.MPopRequest{Key: []byte("h"), Keys: keys})
	log.Printf("mpopResponse: %+v, err: %+v", mpopResponse, err)

	// 发起 mexist 请求
	mexistResponse, err := c.MExist(context.Background(),
		&pb.MExistRequest{Key: []byte("h"),
			Keys: [][]byte{[]byte("k1"), []byte("k2"), []byte("k3000"), []byte("k4000"), []byte("k5")}})
	log.Printf("mexistResponse: %+v, err: %+v", mexistResponse, err)

	// 发起 mcount 请求
	mcountResponse, err := c.MCount(context.Background(),
		&pb.MCountRequest{Key: []byte("h")})
	log.Printf("mcountResponse: %+v, err: %+v", mcountResponse, err)

	// 发起 sdel 请求
	mdelResponse, err := c.MDel(context.Background(),
		&pb.MDelRequest{Key: []byte("h")})
	log.Printf("mdelResponse: %+v, err: %+v", mdelResponse, err)

	mmembersResponse, _ = c.MMembers(context.Background(),
		&pb.MMembersRequest{Key: []byte("h")})
	log.Printf("mmembersResponse: %+v, err: %+v", mmembersResponse, err)
}
