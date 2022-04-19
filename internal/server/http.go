package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var httpLogger = util.GetLogger("http")

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
		httpLogger.Fatalf("failed to register: %+v", err)
	}
	return &HttpServer{mux: mux}
}

func (httpServer *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.RequestURI == "/metrics" {
		promhttp.Handler().ServeHTTP(writer, request)
	} else if strings.HasPrefix(request.RequestURI, "/join") {
		nodeIdStr := request.URL.Query()["nodeId"][0]
		var nodeId uint64
		nodeId, err := strconv.ParseUint(nodeIdStr, 10, 64)
		if err != nil {
			writer.WriteHeader(401)
			_, _ = writer.Write([]byte("failed"))
			return
		}

		address := request.URL.Query()["address"][0]
		if err := store.HandleJoin(nodeId, address); err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("failed"))
			return
		}
		writer.WriteHeader(200)
		_, _ = writer.Write([]byte("ok"))
	} else if strings.HasPrefix(request.RequestURI, "/delete") {
		nodeIdStr := request.URL.Query()["nodeId"][0]
		var nodeId uint64
		nodeId, err := strconv.ParseUint(nodeIdStr, 10, 64)
		if err != nil {
			writer.WriteHeader(401)
			_, _ = writer.Write([]byte("failed"))
			return
		}
		if err := store.HandleDelete(nodeId); err != nil {
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte("failed"))
			return
		}
		writer.WriteHeader(200)
		_, _ = writer.Write([]byte("ok"))
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
		httpLogger.Printf("failed to serve: %+v", err)
	}
}

func (httpServer *HttpServer) Stop() {
	if httpServer.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := httpServer.server.Shutdown(ctx); err != nil {
			httpLogger.Printf("shutdown http error: %+v", err)
		}
		httpLogger.Println("stop http server finished")
	}
}
