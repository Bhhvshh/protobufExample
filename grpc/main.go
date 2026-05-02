package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "example.com/grpcdemo"
	"google.golang.org/grpc"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[pb.ChatService_ChatServer]struct{}
}

func newServer() *chatServer {
	return &chatServer{
		clients: make(map[pb.ChatService_ChatServer]struct{}),
	}
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	s.mu.Unlock()

	log.Println("New client connected")

	for {
		// Receive a message from the client
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		// Broadcast to all connected clients
		s.broadcast(msg)
	}

	// Remove client when disconnected
	s.mu.Lock()
	delete(s.clients, stream)
	s.mu.Unlock()

	log.Println("Client disconnected")
	return nil
}

func (s *chatServer) broadcast(msg *pb.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		err := client.Send(msg)
		if err != nil {
			log.Printf("Error sending to client: %v", err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, newServer())

	fmt.Println("🚀 Chat server started at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
