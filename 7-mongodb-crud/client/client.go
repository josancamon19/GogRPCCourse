package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
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
	fmt.Printf("Blog created: %v\n", res)

	readBlogRequest := &blog.ReadBlogRequest{BlogId: res.GetBlog().GetId()}
	readBlogResponse, err := server.ReadBlog(context.Background(), readBlogRequest)
	if err != nil {
		log.Fatalf("Error retrieving blog: %v\n", err)
	}
	fmt.Printf("Read blog response: %v\n", readBlogResponse)

	listResponse, err := server.ListBlogs(context.Background(), &blog.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error reading list of blogs: %v\n", err)
	}

	for {
		item, err := listResponse.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error receiving: %v\n", err)
		}
		fmt.Println(item.GetBlog().String())
	}
}
