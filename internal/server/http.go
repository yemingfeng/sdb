package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/store"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HttpServer struct {
	mux    *runtime.ServeMux
	server *http.Server
}

func NewHttpServer() *HttpServer {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterSDBHandlerFromEndpoint(context.Background(),
		mux, ":"+strconv.Itoa(conf.Conf.Server.GRPCPort), opts)
	if err != nil {
		log.Fatalf("failed to register: %+v", err)
	}
	return &HttpServer{mux: mux}
}

func (httpServer *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.RequestURI == "/metrics" {
		promhttp.Handler().ServeHTTP(writer, request)
	} else if strings.HasPrefix(request.RequestURI, "/join") {
		address := request.URL.Query()["address"][0]
		nodeId := request.URL.Query()["nodeId"][0]
		if err := store.HandleJoin(address, nodeId); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(request.RequestURI, "/v1") {
		httpServer.mux.ServeHTTP(writer, request)
	} else {
		writer.WriteHeader(502)
	}
}

func (httpServer *HttpServer) Start() {
	server := &http.Server{Addr: ":" + strconv.Itoa(conf.Conf.Server.HttpPort), Handler: httpServer}
	httpServer.server = server

	if err := server.ListenAndServe(); err != nil {
		log.Printf("failed to serve: %+v", err)
	}
}

func (httpServer *HttpServer) Stop() {
	if httpServer.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := httpServer.server.Shutdown(ctx); err != nil {
			log.Printf("shutdown http error: %+v", err)
		}
		log.Println("stop http server finished")
	}
}
