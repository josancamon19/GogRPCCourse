package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
	"udemyGRPC/4-server-streaming/greetpb"
)

type server struct{}

func (s *server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := request.GetGreeting().FirstName
	for i := 0; i < 10; i++ {
		response := greetpb.GreetManyTimesResponse{Result: "Hello " + firstName + " -> " + strconv.Itoa(i)}
		err := stream.Send(&response)
		if err != nil {
			fmt.Printf("Error streaming response: %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
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
