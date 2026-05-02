package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "example.com/grpcdemo/chat"
	"google.golang.org/grpc"
)

func main() {
	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// Ask username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Goroutine: receive messages
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving: %v", err)
				return
			}
			fmt.Printf("[%s]: %s\n", in.User, in.Message)
		}
	}()

	// Main loop: send messages
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		err := stream.Send(&pb.ChatMessage{
			User:    username,
			Message: text,
		})
		if err != nil {
			log.Printf("Error sending: %v", err)
			break
		}
	}
}
