import pika
import json
from datetime import datetime

def callback(ch, method, properties, body):
    message = json.loads(body)
    print(f"Received message: {message}")
    print(f"Name: {message['name']}")
    print(f"Greeting: {message['message']}")
    print(f"Timestamp: {message['timestamp']}")
    print("-" * 50)

def main():
    # RabbitMQ connection
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host='localhost')
    )
    channel = connection.channel()
    
    # Declare the queue
    channel.queue_declare(queue='greetings')
    
    # Set up the consumer
    channel.basic_consume(
        queue='greetings',
        on_message_callback=callback,
        auto_ack=True
    )
    
    print("Starting to consume messages from RabbitMQ...")
    print("Press Ctrl+C to exit")
    
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print("\nStopping message consumption...")
    finally:
        connection.close()

if __name__ == '__main__':
    main()