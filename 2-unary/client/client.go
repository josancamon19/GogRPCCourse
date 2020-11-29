package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"udemyGRPC/2-unary/greetpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	server := greetpb.NewGreetServiceClient(conn)

	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Joan",
		LastName:  "Cabezas",
	}}

	res, err := server.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet RPC %v", err)
	}

	fmt.Println(res.Result)
}
