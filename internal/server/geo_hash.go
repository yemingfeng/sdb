package server

import (
	"github.com/yemingfeng/sdb/internal/service"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
)

type GeoHashServer struct {
	pb.UnimplementedSDBServer
}

func (server *GeoHashServer) GHCreate(_ context.Context, request *pb.GHCreateRequest) (*pb.GHCreateResponse, error) {
	err := service.GHCreate(request.Key, request.Precision)
	return &pb.GHCreateResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHDel(_ context.Context, request *pb.GHDelRequest) (*pb.GHDelResponse, error) {
	err := service.GHDel(request.Key)
	return &pb.GHDelResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHAdd(_ context.Context, request *pb.GHAddRequest) (*pb.GHAddResponse, error) {
	err := service.GHAdd(request.Key, request.Points)
	return &pb.GHAddResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHPop(_ context.Context, request *pb.GHPopRequest) (*pb.GHPopResponse, error) {
	err := service.GHPop(request.Key, request.Ids)
	return &pb.GHPopResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHGetBoxes(_ context.Context, request *pb.GHGetBoxesRequest) (*pb.GHGetBoxesResponse, error) {
	res, err := service.GHGetBoxes(request.Key, request.Latitude, request.Longitude)
	return &pb.GHGetBoxesResponse{Points: res}, err
}

func (server *GeoHashServer) GHGetNeighbors(_ context.Context, request *pb.GHGetNeighborsRequest) (*pb.GHGetNeighborsResponse, error) {
	res, err := service.GHGetNeighbors(request.Key, request.Latitude, request.Longitude)
	return &pb.GHGetNeighborsResponse{Points: res}, err
}

func (server *GeoHashServer) GHCount(_ context.Context, request *pb.GHCountRequest) (*pb.GHCountResponse, error) {
	res, err := service.GHCount(request.Key)
	return &pb.GHCountResponse{Count: res}, err
}

func (server *GeoHashServer) GHMembers(_ context.Context, request *pb.GHMembersRequest) (*pb.GHMembersResponse, error) {
	res, err := service.GHMembers(request.Key)
	return &pb.GHMembersResponse{Points: res}, err
}
