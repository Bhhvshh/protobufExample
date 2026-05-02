package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/user/lpw/proto/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedMessageServiceServer
    collection *mongo.Collection
}

func (s *server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
    // Create MongoDB document
    document := bson.M{
        "messageId":  msg.Id,
        "content":   msg.Content,
        "timestamp": time.Now().Unix(),
    }

    // Insert into MongoDB
    _, err := s.collection.InsertOne(ctx, document)
    if err != nil {
        return &pb.Response{
            Status:  "error",
            Message: "Failed to store message",
        }, err
    }

    return &pb.Response{
        Status:  "success",
        Message: "Message stored successfully",
    }, nil
}

func main() {
    // Connect to MongoDB
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collection := client.Database("messagedb").Collection("messages")

    // Create gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterMessageServiceServer(s, &server{collection: collection})

    log.Printf("Server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}