package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"udemyGRPC/3-exercise-unary/protos"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	server := protos.NewSumAPIClient(conn)
	request := protos.APIRequest{
		Integer1: 10,
		Integer2: 3,
	}
	response, err := server.Sum(context.Background(), &request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Result)
}
