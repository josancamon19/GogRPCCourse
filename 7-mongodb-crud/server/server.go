package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
}

type server struct{}

func (s *server) ListBlogs(_ *blog.ListBlogRequest, stream blog.BlogService_ListBlogsServer) error {
	blogsCursor, err := mongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Error retrieving list from mongoDB: %v", err))
	}
	defer blogsCursor.Close(context.Background())
	for blogsCursor.Next(context.Background()) {
		data := blogItem{}
		err := blogsCursor.Decode(&data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error decoding data from mongoDB: %v", err))
		}
		fmt.Println(data)

		_ = stream.Send(&blog.ListBlogResponse{Blog: &blog.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		}})
	}

	if err := blogsCursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknwon internal error: %v", err))
	}
	return nil
}

func (s *server) ReadBlog(ctx context.Context, request *blog.ReadBlogRequest) (*blog.ReadBlogResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.GetBlogId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID"))
	}
	result := mongoCollection.FindOne(context.Background(), bson.M{"_id": oid})
	var blogFound blogItem
	if err := result.Decode(&blogFound); err != nil {
		return nil, status.Errorf(
			codes.NotFound, fmt.Sprintf("Blog with id %s not found: %v", oid.Hex(), err))
	}
	return &blog.ReadBlogResponse{Blog: &blog.Blog{
		Id:       blogFound.ID.Hex(),
		AuthorId: blogFound.AuthorID,
		Title:    blogFound.Title,
		Content:  blogFound.Content,
	}}, nil
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
