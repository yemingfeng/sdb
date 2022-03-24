package test

import (
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestRaft(t *testing.T) {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		log.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// 连接服务器
	c := pb.NewSDBClient(conn)

	s := semaphore.NewWeighted(200)

	randBytes := func() []byte {
		return []byte("hello" + strconv.Itoa(rand.Int()%10000))
	}

	set := func(key, value []byte) {
		_, err := c.Set(context.Background(), &pb.SetRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("%+v, key = %s, value = %s", err, key, value)
		}
	}

	get := func(key []byte) {
		_, err := c.Get(context.Background(), &pb.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("%+v, key = %s", err, key)
		}
	}

	for i := 0; i < 1000; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				_ = s.Acquire(context.Background(), 1)
				set(randBytes(), randBytes())
				wg.Done()
				s.Release(1)
			}()

			wg.Add(1)
			go func() {
				_ = s.Acquire(context.Background(), 1)
				get(randBytes())
				wg.Done()
				s.Release(1)
			}()
		}
		wg.Wait()
	}
}
