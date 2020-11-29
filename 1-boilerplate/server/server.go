package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"udemyGRPC/1-boilerplate/greetpb"
)

type server struct {

}
func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
