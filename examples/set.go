package main

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var setLogger = util.GetLogger("set")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		setLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	c := pb.NewSDBClient(conn)
	values := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		values[i] = []byte("h" + fmt.Sprint(i))
		values[i] = []byte("h" + fmt.Sprint(i))
	}
	spushResponse, err := c.SPush(context.Background(),
		&pb.SPushRequest{Key: []byte("h"), Values: values})
	setLogger.Printf("spushResponse: %+v, err: %+v", spushResponse, err)

	spushResponse, err = c.SPush(context.Background(),
		&pb.SPushRequest{Key: []byte("h1"), Values: values})
	setLogger.Printf("spushResponse: %+v, err: %+v", spushResponse, err)

	smembersResponse, _ := c.SMembers(context.Background(),
		&pb.SMembersRequest{Key: []byte("h")})
	setLogger.Printf("smembersResponse: %+v, err: %+v", smembersResponse, err)

	smembersResponse, _ = c.SMembers(context.Background(),
		&pb.SMembersRequest{Key: []byte("h1")})
	setLogger.Printf("smembersResponse: %+v, err: %+v", smembersResponse, err)

	values = make([][]byte, 50)
	for i := 0; i < 50; i++ {
		values[i] = []byte("h" + fmt.Sprint(i*2))
	}
	spopResponse, err := c.SPop(context.Background(),
		&pb.SPopRequest{Key: []byte("h"), Values: values})
	setLogger.Printf("spopResponse: %+v, err: %+v", spopResponse, err)

	smembersResponse, _ = c.SMembers(context.Background(),
		&pb.SMembersRequest{Key: []byte("h")})
	setLogger.Printf("smembersResponse: %+v, err: %+v", smembersResponse, err)

	sexistResponse, err := c.SExist(context.Background(),
		&pb.SExistRequest{Key: []byte("h"),
			Values: [][]byte{[]byte("h1"), []byte("h2"), []byte("h3"), []byte("h4"), []byte("h5")}})
	setLogger.Printf("sexistResponse: %+v, err: %+v", sexistResponse, err)

	scountResponse, err := c.SCount(context.Background(),
		&pb.SCountRequest{Key: []byte("h")})
	setLogger.Printf("scountResponse: %+v, err: %+v", scountResponse, err)

	//发起 sdel 请求
	sdelResponse, err := c.SDel(context.Background(),
		&pb.SDelRequest{Key: []byte("h")})
	setLogger.Printf("sdelResponse: %+v, err: %+v", sdelResponse, err)

	smembersResponse, _ = c.SMembers(context.Background(),
		&pb.SMembersRequest{Key: []byte("h")})
	setLogger.Printf("smembersResponse: %+v, err: %+v", smembersResponse, err)
}
