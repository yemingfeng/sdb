package main

import (
	"github.com/yemingfeng/sdb/internal/server"
	"github.com/yemingfeng/sdb/internal/store"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	store.StartRaft()

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
	log.Printf("os signal: %+v", s)

	store.Stop()
	sdbGrpcServer.Stop()
	httpServer.Stop()
}
