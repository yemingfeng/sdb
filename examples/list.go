package main

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var listLogger = util.GetLogger("list")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		listLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	length := 10
	values := make([][]byte, length)
	for i := 0; i < length; i++ {
		values[i] = []byte("h" + fmt.Sprint(i))
	}
	lpushResponse, _ := c.LRPush(context.Background(),
		&pb.LRPushRequest{Key: []byte("h"), Values: values})
	listLogger.Printf("lpushResponse: %+v, err: %+v", lpushResponse, err)

	values = make([][]byte, length)
	for i := 0; i < length; i++ {
		values[i] = []byte("h" + fmt.Sprint((i+1)*100))
	}
	llpushResponse, _ := c.LLPush(context.Background(),
		&pb.LLPushRequest{Key: []byte("h"), Values: values})
	listLogger.Printf("llpushResponse: %+v, err: %+v", llpushResponse, err)
	llpushResponse, _ = c.LLPush(context.Background(),
		&pb.LLPushRequest{Key: []byte("h"), Values: values})
	listLogger.Printf("llpushResponse: %+v, err: %+v", llpushResponse, err)

	lmembersResponse, _ := c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h")})
	listLogger.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)

	length = 50
	values = make([][]byte, length)
	for i := 0; i < length; i++ {
		values[i] = []byte("h" + fmt.Sprint(i*2))
	}
	lpopResponse, err := c.LPop(context.Background(),
		&pb.LPopRequest{Key: []byte("h"), Values: values})
	listLogger.Printf("lpopResponse: %+v, err: %+v", lpopResponse, err)

	lmembersResponse, _ = c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h")})
	listLogger.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)

	// 反向检索
	lrangeResponse, err := c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: -1, Limit: 10})
	listLogger.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: -3, Limit: 3})
	listLogger.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: -4, Limit: 3})
	listLogger.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	// 正向检索
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: 1, Limit: 3})
	listLogger.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: 2, Limit: 10})
	listLogger.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: 0, Limit: 10})
	listLogger.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)

	lexistResponse, err := c.LExist(context.Background(),
		&pb.LExistRequest{Key: []byte("h"), Values: [][]byte{[]byte("h1"),
			[]byte("h2"), []byte("h3"), []byte("h4"), []byte("h5")}})
	listLogger.Printf("lexistResponse: %+v, err: %+v", lexistResponse, err)

	lcountResponse, err := c.LCount(context.Background(),
		&pb.LCountRequest{Key: []byte("h")})
	listLogger.Printf("lcountResponse: %+v, err: %+v", lcountResponse, err)

	//ldelResponse, err := c.LDel(context.Background(),
	//	&pb.LDelRequest{Key: []byte("h")})
	//listLogger.Printf("ldelResponse: %+v, err: %+v", ldelResponse, err)

	lmembersResponse, _ = c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h")})
	listLogger.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)

	llpushResponse, err = c.LLPush(context.Background(),
		&pb.LLPushRequest{Key: []byte("h2"), Values: [][]byte{[]byte("h1"), []byte("h2")}})
	listLogger.Printf("llpushResponse: %+v, err: %+v", llpushResponse, err)

	lmembersResponse, err = c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h2")})
	listLogger.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)

	lpopResponse, err = c.LPop(context.Background(),
		&pb.LPopRequest{Key: []byte("h2"), Values: [][]byte{[]byte("h1"), []byte("h2")}})
	listLogger.Printf("lpopResponse: %+v, err: %+v", lpopResponse, err)
}
