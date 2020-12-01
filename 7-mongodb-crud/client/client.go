package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	blog "udemyGRPC/7-mongodb-crud/protos"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	server := blog.NewBlogServiceClient(conn)

	request := blog.CreateBlogRequest{Blog: &blog.Blog{
		AuthorId: "josancamon19",
		Title:    "Medium post",
		Content:  "Medium post content",
	}}
	res, err := server.CreateBlog(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error calling Greet RPC %v", err)
	}
	fmt.Printf("Blog created: %v", res)
}
