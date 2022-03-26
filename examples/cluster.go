package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var clusterLogger = util.GetLogger("cluster")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		clusterLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// 连接服务器
	c := pb.NewSDBClient(conn)
	cinfoResponse, err := c.CInfo(context.Background(), &pb.CInfoRequest{})
	clusterLogger.Printf("cinfoResponse: %+v, err: %+v", cinfoResponse, err)
}
