package main

import (
	context "context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"udemyGRPC/3-exercise-unary/protos"
)

type server struct{}

func (s *server) Sum(ctx context.Context, request *protos.APIRequest) (*protos.APIResponse, error) {
	return &protos.APIResponse{Result: request.GetInteger1() + request.GetInteger2()}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	protos.RegisterSumAPIServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
