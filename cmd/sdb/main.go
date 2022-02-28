package main

import (
	"github.com/yemingfeng/sdb/internal/server"
)

func main() {
	httpServer := server.NewHttpServer()
	go func() {
		httpServer.Start()
	}()

	sdbGrpcServer := server.NewSDBGrpcServer()
	sdbGrpcServer.Start()
}
