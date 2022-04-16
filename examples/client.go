package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var clientLogger = util.GetLogger("client")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		clientLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	setResponse, err := c.Set(context.Background(),
		&pb.SetRequest{Key: []byte("hello"), Value: []byte("world")})
	clientLogger.Printf("setResponse: %+v, err: %+v", setResponse, err)
	getResponse, err := c.Get(context.Background(),
		&pb.GetRequest{Key: []byte("hello")})
	clientLogger.Printf("getResponse: %+v, err: %+v", getResponse, err)
}
