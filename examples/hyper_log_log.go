package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var hyperLogLogLogger = util.GetLogger("hyper_log_log")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		hyperLogLogLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	hllCreateResponse, err := c.HLLCreate(context.Background(),
		&pb.HLLCreateRequest{Key: []byte("hello")})
	hyperLogLogLogger.Printf("hllCreateResponse: %+v, err: %+v", hllCreateResponse, err)
	hllAddResponse, err := c.HLLAdd(context.Background(),
		&pb.HLLAddRequest{Key: []byte("hello"),
			Values: [][]byte{[]byte("aaa"), []byte("bbb"), []byte("ccc"),
				[]byte("ddd"), []byte("aaa"), []byte("eee"), []byte("bbb")}})
	hyperLogLogLogger.Printf("hllAddResponse: %+v, err: %+v", hllAddResponse, err)
	hllCountResponse, err := c.HLLCount(context.Background(),
		&pb.HLLCountRequest{Key: []byte("hello")})
	hyperLogLogLogger.Printf("hllCountResponse: %+v, err: %+v", hllCountResponse, err)
	hllDelResponse, err := c.HLLDel(context.Background(),
		&pb.HLLDelRequest{Key: []byte("hello")})
	hyperLogLogLogger.Printf("hllDelResponse: %+v, err: %+v", hllDelResponse, err)
	hllCountResponse, err = c.HLLCount(context.Background(),
		&pb.HLLCountRequest{Key: []byte("hello")})
	hyperLogLogLogger.Printf("hllCountResponse: %+v, err: %+v", hllCountResponse, err)
}
