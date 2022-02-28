package server

import (
	"context"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type SDBGrpcServer struct {
	grpcServer *grpc.Server
	StringServer
	ListServer
	SetServer
	SortedSetServer
	BloomFilterServer
	HyperLogLogServer
	BitsetServer
	MapServer
	GeoHashServer
	PubSubServer
}

func NewSDBGrpcServer() *SDBGrpcServer {
	grpcprometheus.EnableHandlingTimeHistogram()

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				err := handler(srv, ss)
				if err != nil {
					log.Println(err)
					fmt.Errorf("error: %v", err)
				}
				return err
			},
			grpcprometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				resp, err = handler(ctx, req)
				if err != nil {
					log.Println(err)
					fmt.Errorf("error: %v", err)
				}
				return resp, err
			},
			grpcmiddleware.ChainUnaryServer(
				ratelimit.UnaryServerInterceptor(CreateRateLimit(conf.Conf.Server.Rate))),
			grpcprometheus.UnaryServerInterceptor,
			grpcmiddleware.ChainUnaryServer(func(ctx context.Context, req interface{},
				info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				start := time.Now()

				defer func() {
					cost := time.Since(start)
					if cost.Milliseconds() > conf.Conf.Server.SlowQueryThreshold {
						log.Printf("Slow query: %s, cost: %d", req, cost.Milliseconds())
					}
				}()

				return handler(ctx, req)
			}),
		)),
	)
	sdbGrpcServer := SDBGrpcServer{grpcServer: grpcServer}
	pb.RegisterSDBServer(grpcServer, &sdbGrpcServer)
	return &sdbGrpcServer
}

func (sdbGrpcServer *SDBGrpcServer) Start() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Conf.Server.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
	}
	if err := sdbGrpcServer.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
	log.Printf("serve: %d", conf.Conf.Server.GRPCPort)
}
