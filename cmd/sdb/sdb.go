package main

import (
	"github.com/yemingfeng/sdb/internal/server"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var sdbLogger = util.GetLogger("sdb")

func main() {
	store.StartCluster()

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
	sdbLogger.Printf("receive os signal: %+v", s)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(4)

	go func() {
		store.Stop()
		defer waitGroup.Done()
	}()
	go func() {
		store.StopCluster()
		defer waitGroup.Done()
	}()
	go func() {
		sdbGrpcServer.Stop()
		defer waitGroup.Done()
	}()
	go func() {
		httpServer.Stop()
		waitGroup.Done()
	}()

	waitGroup.Wait()
}
