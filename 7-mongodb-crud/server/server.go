package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"
	blog "udemyGRPC/7-mongodb-crud/protos"
)

var mongoClient *mongo.Client
var mongoCollection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `jbon:"title"`
}

func initMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://josancamon19:123@cluster0.kigjv.mongodb.net/cluster0?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection = mongoClient.Database("cluster0").Collection("blog")
	fmt.Println(mongoCollection)
}

type server struct {
}

func (s *server) CreateBlog(ctx context.Context, request *blog.CreateBlogRequest) (*blog.CreateBlogResponse, error) {
	blogData := request.GetBlog()
	data := blogItem{
		AuthorID: blogData.AuthorId,
		Content:  blogData.Content,
		Title:    blogData.Title,
	}
	res, err := mongoCollection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID"))
	}

	return &blog.CreateBlogResponse{Blog: &blog.Blog{
		Id:       oid.Hex(),
		AuthorId: blogData.GetAuthorId(),
		Title:    blogData.GetTitle(),
		Content:  blogData.GetContent(),
	}}, nil

}

func main() {
	initMongo()
	lis, _ := net.Listen("tcp", "0.0.0.0:50051")
	grpcServer := grpc.NewServer()
	blog.RegisterBlogServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
