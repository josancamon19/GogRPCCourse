package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
	"time"
	"udemyGRPC/6-bidirectional-streaming/greetpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	server := greetpb.NewGreetServiceClient(conn)

	requests := []*greetpb.GreetEveryoneRequest{
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

	stream, err := server.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error calling GreetEveryone RPC %v", err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer stream.CloseSend()
		for _, request := range requests {
			fmt.Printf("Sending request %s\n", request)
			err = stream.Send(request)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
		waitGroup.Done()
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error receiving: %v", err)
			}
			fmt.Printf("Receving response %s\n", response.GetResult())
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
}
