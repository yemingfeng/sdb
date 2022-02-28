package server

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type ListServer struct {
	pb.UnimplementedSDBServer
}

func (server *ListServer) LRPush(_ context.Context, request *pb.LRPushRequest) (*pb.LRPushResponse, error) {
	err := service.LRPush(request.Key, request.Values)
	return &pb.LRPushResponse{Success: err == nil}, err
}

func (server *ListServer) LLPush(_ context.Context, request *pb.LLPushRequest) (*pb.LLPushResponse, error) {
	err := service.LLPush(request.Key, request.Values)
	return &pb.LLPushResponse{Success: err == nil}, err
}

func (server *ListServer) LPop(_ context.Context, request *pb.LPopRequest) (*pb.LPopResponse, error) {
	err := service.LPop(request.Key, request.Values)
	return &pb.LPopResponse{Success: err == nil}, err
}

func (server *ListServer) LRange(_ context.Context, request *pb.LRangeRequest) (*pb.LRangeResponse, error) {
	res, err := service.LRange(request.Key, request.Offset, request.Limit)
	return &pb.LRangeResponse{Values: res}, err
}

func (server *ListServer) LExist(_ context.Context, request *pb.LExistRequest) (*pb.LExistResponse, error) {
	res, err := service.LExist(request.Key, request.Values)
	return &pb.LExistResponse{Exists: res}, err
}

func (server *ListServer) LDel(_ context.Context, request *pb.LDelRequest) (*pb.LDelResponse, error) {
	err := service.LDel(request.Key)
	return &pb.LDelResponse{Success: err == nil}, err
}

func (server *ListServer) LCount(_ context.Context, request *pb.LCountRequest) (*pb.LCountResponse, error) {
	res, err := service.LCount(request.Key)
	return &pb.LCountResponse{Count: res}, err
}

func (server *ListServer) LMembers(_ context.Context, request *pb.LMembersRequest) (*pb.LMembersResponse, error) {
	res, err := service.LMembers(request.Key)
	return &pb.LMembersResponse{Values: res}, err
}
