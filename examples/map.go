package main

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var mapLogger = util.GetLogger("map")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		mapLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	pairs := make([]*pb.Pair, 100)
	for i := 0; i < 100; i++ {
		pairs[i] = &pb.Pair{Key: []byte("k" + fmt.Sprint(i)), Value: []byte("v" + fmt.Sprint(i+1))}
	}
	mpushResponse, err := c.MPush(context.Background(),
		&pb.MPushRequest{Key: []byte("h"), Pairs: pairs})
	mapLogger.Printf("mpushResponse: %+v, err: %+v", mpushResponse, err)

	mmembersResponse, _ := c.MMembers(context.Background(),
		&pb.MMembersRequest{Key: []byte("h")})
	mapLogger.Printf("mmembersResponse: %+v, err: %+v", mmembersResponse, err)

	keys := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		keys[i] = []byte("k" + fmt.Sprint(i*2))
	}
	mpopResponse, err := c.MPop(context.Background(),
		&pb.MPopRequest{Key: []byte("h"), Keys: keys})
	mapLogger.Printf("mpopResponse: %+v, err: %+v", mpopResponse, err)

	mexistResponse, err := c.MExist(context.Background(),
		&pb.MExistRequest{Key: []byte("h"),
			Keys: [][]byte{[]byte("k1"), []byte("k2"), []byte("k3000"), []byte("k4000"), []byte("k5")}})
	mapLogger.Printf("mexistResponse: %+v, err: %+v", mexistResponse, err)

	mcountResponse, err := c.MCount(context.Background(),
		&pb.MCountRequest{Key: []byte("h")})
	mapLogger.Printf("mcountResponse: %+v, err: %+v", mcountResponse, err)

	//mdelRespo/**/nse, err := c.MDel(context.Background(),
	//	&pb.MDelRequest{Key: []byte("h")})
	//mapLogger.Printf("mdelResponse: %+v, err: %+v", mdelResponse, err)

	mmembersResponse, _ = c.MMembers(context.Background(),
		&pb.MMembersRequest{Key: []byte("h")})
	mapLogger.Printf("mmembersResponse: %+v, err: %+v", mmembersResponse, err)

	mpushResponse, err = c.MPush(context.Background(),
		&pb.MPushRequest{Key: []byte("h2"),
			Pairs: []*pb.Pair{{Key: []byte("h1"), Value: []byte("h2")},
				{Key: []byte("h3"), Value: []byte("h4")}}})
	mapLogger.Printf("mpushResponse: %+v, err: %+v", mpushResponse, err)

	mmembersResponse, err = c.MMembers(context.Background(),
		&pb.MMembersRequest{Key: []byte("h2")})
	mapLogger.Printf("mmembersResponse: %+v, err: %+v", mmembersResponse, err)

	mpopResponse, err = c.MPop(context.Background(),
		&pb.MPopRequest{Key: []byte("h2"), Keys: [][]byte{[]byte("h1"), []byte("h3")}})
	mapLogger.Printf("mpopResponse: %+v, err: %+v", mpopResponse, err)
}
