package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
	"udemyGRPC/5-client-streaming/greetpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	server := greetpb.NewGreetServiceClient(conn)

	requests := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "Joan",
			LastName:  "Cabezas",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Juan",
			LastName:  "Pepino",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Santiago",
			LastName:  "Monroy",
		}},
	}

	stream, err := server.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error calling LongGreet RPC %v", err)
	}
	for _, request := range requests {
		fmt.Printf("Sending request %s\n", request)
		err = stream.Send(request)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response %v", err)
	}
	fmt.Println(response)
}
