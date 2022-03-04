package cli

import (
	"flag"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/grpc"
	"log"
)

var client pb.SDBClient

func init() {
	server := flag.String("server", "127.0.0.1:10000", "sdb server")
	conn, err := grpc.Dial(*server, grpc.WithInsecure())
	if err != nil {
		log.Printf("faild to connect: %s, %+v", *server, err)
	}
	client = pb.NewSDBClient(conn)
}
