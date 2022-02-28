package main

import (
	"fmt"
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
	// 发起 lpush 请求
	length := 10
	values := make([][]byte, length)
	for i := 0; i < length; i++ {
		values[i] = []byte("h" + fmt.Sprint(i))
	}
	lpushResponse, _ := c.LRPush(context.Background(),
		&pb.LRPushRequest{Key: []byte("h"), Values: values})
	log.Printf("lpushResponse: %+v, err: %+v", lpushResponse, err)

	// 发起 llpush 请求
	values = make([][]byte, length)
	for i := 0; i < length; i++ {
		values[i] = []byte("h" + fmt.Sprint((i+1)*100))
	}
	llpushResponse, _ := c.LLPush(context.Background(),
		&pb.LLPushRequest{Key: []byte("h"), Values: values})
	log.Printf("llpushResponse: %+v, err: %+v", llpushResponse, err)
	llpushResponse, _ = c.LLPush(context.Background(),
		&pb.LLPushRequest{Key: []byte("h"), Values: values})
	log.Printf("llpushResponse: %+v, err: %+v", llpushResponse, err)

	lmembersResponse, _ := c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h")})
	log.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)

	// 发起 lpop 请求
	length = 50
	values = make([][]byte, length)
	for i := 0; i < length; i++ {
		values[i] = []byte("h" + fmt.Sprint(i*2))
	}
	lpopResponse, err := c.LPop(context.Background(),
		&pb.LPopRequest{Key: []byte("h"), Values: values})
	log.Printf("lpopResponse: %+v, err: %+v", lpopResponse, err)

	lmembersResponse, _ = c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h")})
	log.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)

	// 发起 lrange 请求
	// 反向检索
	lrangeResponse, err := c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: -1, Limit: 10})
	log.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: -3, Limit: 3})
	log.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: -4, Limit: 3})
	log.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	// 正向检索
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: 1, Limit: 3})
	log.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: 2, Limit: 10})
	log.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)
	lrangeResponse, err = c.LRange(context.Background(),
		&pb.LRangeRequest{Key: []byte("h"), Offset: 0, Limit: 10})
	log.Printf("lrangeResponse: %+v, err: %+v", lrangeResponse, err)

	// 发起 lexist 请求
	lexistResponse, err := c.LExist(context.Background(),
		&pb.LExistRequest{Key: []byte("h"), Values: [][]byte{[]byte("h1"),
			[]byte("h2"), []byte("h3"), []byte("h4"), []byte("h5")}})
	log.Printf("lexistResponse: %+v, err: %+v", lexistResponse, err)

	// 发起 lcount 请求
	lcountResponse, err := c.LCount(context.Background(),
		&pb.LCountRequest{Key: []byte("h")})
	log.Printf("lcountResponse: %+v, err: %+v", lcountResponse, err)

	// 发起 ldel 请求
	ldelResponse, err := c.LDel(context.Background(),
		&pb.LDelRequest{Key: []byte("h")})
	log.Printf("ldelResponse: %+v, err: %+v", ldelResponse, err)

	lmembersResponse, _ = c.LMembers(context.Background(),
		&pb.LMembersRequest{Key: []byte("h")})
	log.Printf("lmembersResponse: %+v, err: %+v", lmembersResponse, err)
}
