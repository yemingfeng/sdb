package main

import (
	"github.com/yemingfeng/sdb/internal/server"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/util"
	"os"
	"os/signal"
	"syscall"
)

var sdbLogger = util.GetLogger("sdb")

func main() {
	httpServer := server.NewHttpServer()
	go func() {
		httpServer.Start()
	}()

	sdbGrpcServer := server.NewSDBGrpcServer()
	go func() {
		sdbGrpcServer.Start()
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-c
	sdbLogger.Printf("os signal: %+v", s)

	store.Stop()
	sdbGrpcServer.Stop()
	httpServer.Stop()
}
