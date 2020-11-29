package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"udemyGRPC/5-client-streaming/greetpb"
)

type server struct{}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	var message string
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: message})
		}
		if err != nil {
			log.Fatalf("Error reading client stream %v", err)
		}
		greeting := request.GetGreeting()
		message += "Hello " + greeting.FirstName + "!\n"
	}
}

func main() {
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
