package main

import (
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
	// 发起 hllcreate 请求
	hllCreateResponse, err := c.HLLCreate(context.Background(),
		&pb.HLLCreateRequest{Key: []byte("hello")})
	log.Printf("hllCreateResponse: %+v, err: %+v", hllCreateResponse, err)
	// 发起 hlladd 请求
	hllAddResponse, err := c.HLLAdd(context.Background(),
		&pb.HLLAddRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("bbb"), []byte("ccc"),
				[]byte("ddd"), []byte("aaa"), []byte("eee"), []byte("bbb")}})
	log.Printf("hllAddResponse: %+v, err: %+v", hllAddResponse, err)
	// 发起 hllcount 请求
	hllCountResponse, err := c.HLLCount(context.Background(),
		&pb.HLLCountRequest{Key: []byte("hello")})
	log.Printf("hllCountResponse: %+v, err: %+v", hllCountResponse, err)
	// 发起 hlldel 请求
	hllDelResponse, err := c.HLLDel(context.Background(),
		&pb.HLLDelRequest{Key: []byte("hello")})
	log.Printf("hllDelResponse: %+v, err: %+v", hllDelResponse, err)
	// 发起 hllcount 请求
	hllCountResponse, err = c.HLLCount(context.Background(),
		&pb.HLLCountRequest{Key: []byte("hello")})
	log.Printf("hllCountResponse: %+v, err: %+v", hllCountResponse, err)
}
