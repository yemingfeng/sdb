package main

import (
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

	ghCreateResponse, err := c.GHCreate(context.Background(), &pb.GHCreateRequest{Key: []byte("gh1"), Precision: 2})
	log.Printf("ghCreateResponse: %+v, err: %+v", ghCreateResponse, err)

	ghAddResponse, err := c.GHAdd(context.Background(), &pb.GHAddRequest{Key: []byte("gh1"),
		Points: []*pb.Point{
			{Latitude: 11.11, Longitude: 22.11, Id: []byte("p1")},
			{Latitude: 11.22, Longitude: 22.22, Id: []byte("p2")},
			{Latitude: 11.33, Longitude: 22.33, Id: []byte("p3")},
			{Latitude: 11.00, Longitude: 22.00, Id: []byte("p4")},
			{Latitude: 10.90, Longitude: 21.90, Id: []byte("p5")},
			{Latitude: 11.25, Longitude: 22.25, Id: []byte("p6")},
			{Latitude: 11.30, Longitude: 22.30, Id: []byte("p7")},
			{Latitude: 11.10, Longitude: 22.19, Id: []byte("p8")},
			{Latitude: 11.05, Longitude: 22.05, Id: []byte("p9")},
			{Latitude: 11.05, Longitude: 22.25, Id: []byte("p10")},
			{Latitude: 11.10, Longitude: 22.15, Id: []byte("p9")},
			{Latitude: 11.12, Longitude: 22.17, Id: []byte("p90")},
		},
	})
	log.Printf("ghAddResponse: %+v, err: %+v", ghAddResponse, err)

	ghMembersResponse, err := c.GHMembers(context.Background(), &pb.GHMembersRequest{Key: []byte("gh1")})
	log.Printf("ghMembersResponse: %+v, err: %+v", ghMembersResponse, err)
	ghCountResponse, err := c.GHCount(context.Background(), &pb.GHCountRequest{Key: []byte("gh1")})
	log.Printf("ghCountResponse: %+v, err: %+v", ghCountResponse, err)

	ghRemResponse, err := c.GHRem(context.Background(), &pb.GHRemRequest{Key: []byte("gh1"),
		Ids: [][]byte{[]byte("p1"), []byte("p9")},
	})
	log.Printf("ghRemResponse: %+v, err: %+v", ghRemResponse, err)
	ghMembersResponse, err = c.GHMembers(context.Background(), &pb.GHMembersRequest{Key: []byte("gh1")})
	log.Printf("ghMembersResponse: %+v, err: %+v", ghMembersResponse, err)
	ghCountResponse, err = c.GHCount(context.Background(), &pb.GHCountRequest{Key: []byte("gh1")})
	log.Printf("ghCountResponse: %+v, err: %+v", ghCountResponse, err)

	getBoxesResponse, err := c.GHGetBoxes(context.Background(), &pb.GHGetBoxesRequest{Key: []byte("gh1"),
		Latitude: 11.10, Longitude: 22.11})
	log.Printf("getBoxesResponse: %+v, err: %+v", getBoxesResponse, err)

	getNeighborsResponse, err := c.GHGetNeighbors(context.Background(), &pb.GHGetNeighborsRequest{Key: []byte("gh1"),
		Latitude: 11.10, Longitude: 11.12})
	log.Printf("getNeighborsResponse: %+v, err: %+v", getNeighborsResponse, err)

	ghDelResponse, err := c.GHDel(context.Background(), &pb.GHDelRequest{Key: []byte("gh1")})
	log.Printf("ghDelResponse: %+v, err: %+v", ghDelResponse, err)

	ghMembersResponse, err = c.GHMembers(context.Background(), &pb.GHMembersRequest{Key: []byte("gh1")})
	log.Printf("ghMembersResponse: %+v, err: %+v", ghMembersResponse, err)
}
