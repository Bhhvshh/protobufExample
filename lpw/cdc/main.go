package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
    MessageId string `bson:"messageId"`
    Content   string `bson:"content"`
    Timestamp int64  `bson:"timestamp"`
}

func main() {
    // Connect to MongoDB
    ctx := context.Background()
    mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer mongoClient.Disconnect(ctx)

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

    // Watch MongoDB changes
    collection := mongoClient.Database("messagedb").Collection("messages")
    streamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)
    stream, err := collection.Watch(ctx, mongo.Pipeline{}, streamOptions)
    if err != nil {
        log.Fatalf("Failed to create change stream: %v", err)
    }
    defer stream.Close(ctx)

    log.Println("CDC service started. Watching for changes...")

    for stream.Next(ctx) {
        var changeEvent bson.M
        if err := stream.Decode(&changeEvent); err != nil {
            log.Printf("Error decoding change event: %v", err)
            continue
        }

        fullDocument, ok := changeEvent["fullDocument"].(bson.M)
        if !ok {
            log.Printf("Error: fullDocument not found in change event")
            continue
        }

        message := Message{
            MessageId: fullDocument["messageId"].(string),
            Content:   fullDocument["content"].(string),
            Timestamp: fullDocument["timestamp"].(int64),
        }

        // Publish to RabbitMQ
        messageBytes, err := json.Marshal(message)
        if err != nil {
            log.Printf("Error marshaling message: %v", err)
            continue
        }

        err = ch.Publish(
            "",
            q.Name,
            false,
            false,
            amqp.Publishing{
                ContentType: "application/json",
                Body:        messageBytes,
            },
        )
        if err != nil {
            log.Printf("Error publishing message: %v", err)
            continue
        }

        log.Printf("Published message to queue: %s", message.Content)
    }
}