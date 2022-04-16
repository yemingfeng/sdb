package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var sortedSetLogger = util.GetLogger("sorted_set")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		sortedSetLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	zpushResponse, err := c.ZPush(context.Background(),
		&pb.ZPushRequest{Key: []byte("h"),
			Tuples: []*pb.Tuple{
				{Value: []byte("aaa"), Score: 1.0},
				{Value: []byte("ddd"), Score: 0.8},
				{Value: []byte("bbb"), Score: 1.1},
				{Value: []byte("ccc"), Score: 0.9},
				{Value: []byte("eee"), Score: 0.7},
				{Value: []byte("aaa"), Score: 1.23},
			}})
	sortedSetLogger.Printf("zpushResponse: %+v, err: %+v", zpushResponse, err)
	zmembersResponse, _ := c.ZMembers(context.Background(),
		&pb.ZMembersRequest{Key: []byte("h")})
	sortedSetLogger.Printf("zmembersResponse: %+v, err: %+v", zmembersResponse, err)
	zrangeResponse, err := c.ZRange(context.Background(),
		&pb.ZRangeRequest{Key: []byte("h"), Offset: 1, Limit: 100})
	sortedSetLogger.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
	zrangeResponse, err = c.ZRange(context.Background(),
		&pb.ZRangeRequest{Key: []byte("h"), Offset: -1, Limit: 100})
	sortedSetLogger.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
	zpopResponse, err := c.ZPop(context.Background(),
		&pb.ZPopRequest{Key: []byte("h"), Values: [][]byte{[]byte("aaa"), []byte("bbb")}})
	sortedSetLogger.Printf("zpopResponse: %+v, err: %+v", zpopResponse, err)
	zpopResponse, err = c.ZPop(context.Background(),
		&pb.ZPopRequest{Key: []byte("h"), Values: [][]byte{[]byte("ccc"), []byte("ddd"), []byte("eee")}})
	sortedSetLogger.Printf("zpopResponse: %+v, err: %+v", zpopResponse, err)
	zrangeResponse, err = c.ZRange(context.Background(),
		&pb.ZRangeRequest{Key: []byte("h"), Offset: 0, Limit: 100})
	sortedSetLogger.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
	zexistResponse, err := c.ZExist(context.Background(),
		&pb.ZExistRequest{Key: []byte("h"),
			Values: [][]byte{[]byte("aaa"), []byte("ccc"), []byte("ddd")}})
	sortedSetLogger.Printf("zexistResponse: %+v, err: %+v", zexistResponse, err)
	//zdelResponse, err := c.ZDel(context.Background(), &pb.ZDelRequest{Key: []byte("h")})
	//sortedSetLogger.Printf("zdelResponse: %+v, err: %+v", zdelResponse, err)
	zrangeResponse, err = c.ZRange(context.Background(),
		&pb.ZRangeRequest{Key: []byte("h"), Offset: 0, Limit: 100})
	sortedSetLogger.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
}
