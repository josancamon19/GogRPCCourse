package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"udemyGRPC/6-bidirectional-streaming/greetpb"
)

type server struct{}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error reading client stream %v", err)
			return err
		}

		err = stream.Send(&greetpb.GreetEveryoneResponse{Result: "Hello " + req.GetGreeting().GetFirstName() + "!"})
		if err != nil {
			log.Fatalf("Error sending data to client %v", err)
			return err
		}

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
