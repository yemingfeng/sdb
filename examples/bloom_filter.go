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
	// 发起 bfcreate 请求
	bfCreateResponse, err := c.BFCreate(context.Background(),
		&pb.BFCreateRequest{Key: []byte("hello"), N: 10000, P: 0.05})
	log.Printf("bfCreateResponse: %+v, err: %+v", bfCreateResponse, err)
	// 发起 bfadd 请求
	bfAddResponse, err := c.BFAdd(context.Background(),
		&pb.BFAddRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("bbb"), []byte("ccc"), []byte("ddd")}})
	log.Printf("bfAddResponse: %+v, err: %+v", bfAddResponse, err)
	// 发起 bfexist 请求
	bfExistResponse, err := c.BFExist(context.Background(),
		&pb.BFExistRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("eee"), []byte("ccc")}})
	log.Printf("bfExistResponse: %+v, err: %+v", bfExistResponse, err)
	// 发起 bfdel 请求
	bfDelResponse, err := c.BFDel(context.Background(),
		&pb.BFDelRequest{Key: []byte("hello")})
	log.Printf("bfDelResponse: %+v, err: %+v", bfDelResponse, err)
	// 发起 bfexist 请求
	bfExistResponse, err = c.BFExist(context.Background(),
		&pb.BFExistRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("eee"), []byte("ccc")}})
	log.Printf("bfExistResponse: %+v, err: %+v", bfExistResponse, err)
}
