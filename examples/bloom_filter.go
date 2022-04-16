package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var bloomFilterLogger = util.GetLogger("bloom_filter")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		bloomFilterLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	
	c := pb.NewSDBClient(conn)
	// 发起 bfcreate 请求
	bfCreateResponse, err := c.BFCreate(context.Background(),
		&pb.BFCreateRequest{Key: []byte("hello"), N: 10000, P: 0.05})
	bloomFilterLogger.Printf("bfCreateResponse: %+v, err: %+v", bfCreateResponse, err)
	// 发起 bfadd 请求
	bfAddResponse, err := c.BFAdd(context.Background(),
		&pb.BFAddRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("bbb"), []byte("ccc"), []byte("ddd")}})
	bloomFilterLogger.Printf("bfAddResponse: %+v, err: %+v", bfAddResponse, err)
	// 发起 bfexist 请求
	bfExistResponse, err := c.BFExist(context.Background(),
		&pb.BFExistRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("eee"), []byte("ccc")}})
	bloomFilterLogger.Printf("bfExistResponse: %+v, err: %+v", bfExistResponse, err)
	// 发起 bfdel 请求
	bfDelResponse, err := c.BFDel(context.Background(),
		&pb.BFDelRequest{Key: []byte("hello")})
	bloomFilterLogger.Printf("bfDelResponse: %+v, err: %+v", bfDelResponse, err)
	// 发起 bfexist 请求
	bfExistResponse, err = c.BFExist(context.Background(),
		&pb.BFExistRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("eee"), []byte("ccc")}})
	bloomFilterLogger.Printf("bfExistResponse: %+v, err: %+v", bfExistResponse, err)
}
