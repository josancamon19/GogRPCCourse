package main

import (
	"google.golang.org/grpc"
	"log"
	"udemyGRPC/1-boilerplate/greetpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	greetpb.NewGreetServiceClient(conn)
}
