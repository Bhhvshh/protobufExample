package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
	pb "github.com/user/lpw/proto/message"
	"google.golang.org/grpc"
)

type Message struct {
    MessageId string `json:"messageId"`
    Content   string `json:"content"`
    Timestamp int64  `json:"timestamp"`
}

func main() {
    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    defer conn.Close()

    client := pb.NewMessageServiceClient(conn)

    // Connect to RabbitMQ
    rabbitmqConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer rabbitmqConn.Close()

    ch, err := rabbitmqConn.Channel()
    if err != nil {
        log.Fatalf("Failed to open channel: %v", err)
    }
    defer ch.Close()

    // Declare queue
    q, err := ch.QueueDeclare(
        "messages",
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    msgs, err := ch.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    // Send a test message via gRPC
    go func() {
        time.Sleep(2 * time.Second) // Wait for everything to set up
        testMessage := &pb.Message{
            Id:      "test-1",
            Content: "Hello, this is a test message!",
        }
        
        resp, err := client.SendMessage(context.Background(), testMessage)
        if err != nil {
            log.Printf("Error sending message: %v", err)
        } else {
            log.Printf("Message sent via gRPC. Response: %v", resp)
        }
    }()

    log.Println("Started consuming messages from queue...")
    for msg := range msgs {
        var message Message
        if err := json.Unmarshal(msg.Body, &message); err != nil {
            log.Printf("Error decoding message: %v", err)
            continue
        }

        log.Printf("Received message from queue: ID=%s, Content=%s, Timestamp=%v",
            message.MessageId, message.Content, message.Timestamp)
    }
}