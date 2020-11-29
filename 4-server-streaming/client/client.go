package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"udemyGRPC/4-server-streaming/greetpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	server := greetpb.NewGreetServiceClient(conn)

	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{
		FirstName: "Joan",
		LastName:  "Cabezas",
	}}

	responseStream, err := server.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet RPC %v", err)
	}
	for {
		msg, err := responseStream.Recv()
		if err == io.EOF { // TODO what io.EOF means
			break
		}
		if err != nil {
			log.Fatalf("error reading stream %v", err)
		}

		fmt.Printf("msg received: %s\n", msg.GetResult())
	}
}
