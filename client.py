import grpc
import hello_pb2
import hello_pb2_grpc
import pika
import json
import threading
import time
import queue
import sys
from threading import Event

def handle_rabbitmq_message(message_queue, message_event, body):
    try:
        message = json.loads(body)
        message_queue.put(message)
        message_event.set()  # Signal that a new message is available
    except Exception as e:
        print(f"Error processing message: {e}")

class Client:
    def __init__(self):
        self.message_queue = queue.Queue()
        self.message_event = Event()
        self.is_running = True
        self.setup_grpc()
        self.setup_rabbitmq()

    def setup_grpc(self):
        self.channel = grpc.insecure_channel('localhost:50051')
        self.stub = hello_pb2_grpc.GreeterStub(self.channel)

    def setup_rabbitmq(self):
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host='localhost')
        )
        self.rabbitmq_channel = self.connection.channel()
        self.rabbitmq_channel.queue_declare(queue='greetings')
        
        # Set up consumer with a lambda to pass additional parameters
        self.rabbitmq_channel.basic_consume(
            queue='greetings',
            on_message_callback=lambda ch, method, props, body: (
                handle_rabbitmq_message(self.message_queue, self.message_event, body),
                ch.basic_ack(method.delivery_tag)
            ),
            auto_ack=False
        )

    def process_messages(self):
        while self.is_running:
            # Wait for the message event with a timeout
            if self.message_event.wait(0.1):  # 100ms timeout
                while not self.message_queue.empty():
                    message = self.message_queue.get_nowait()
                    print(f"\nReceived message from queue:")
                    print(f"Name: {message['name']}")
                    print(f"Greeting: {message['message']}")
                    print(f"Timestamp: {message['timestamp']}")
                    print("-" * 50)
                    print("Enter your name (or 'quit' to exit): ", end='', flush=True)
                self.message_event.clear()

    def consume_messages(self):
        try:
            self.rabbitmq_channel.start_consuming()
        except Exception as e:
            print(f"RabbitMQ error: {e}")

    def handle_input(self):
        while self.is_running:
            try:
                name = input("Enter your name (or 'quit' to exit): ")
                if name.lower() == 'quit':
                    self.is_running = False
                    break

                response = self.stub.SayHello(hello_pb2.HelloRequest(name=name))
                print("Server response:", response.message)
            except Exception as e:
                print(f"Error: {e}")

    def run(self):
        # Start RabbitMQ consumer thread
        consumer_thread = threading.Thread(target=self.consume_messages)
        consumer_thread.daemon = True
        consumer_thread.start()

        # Start message processing thread
        processor_thread = threading.Thread(target=self.process_messages)
        processor_thread.daemon = True
        processor_thread.start()

        # Handle input in the main thread
        try:
            self.handle_input()
        except KeyboardInterrupt:
            print("\nStopping client...")
        finally:
            self.stop()

    def stop(self):
        self.is_running = False
        if hasattr(self, 'connection') and self.connection:
            self.connection.close()
        if hasattr(self, 'channel') and self.channel:
            self.channel.close()

def main():
    client = Client()
    client.run()

if __name__ == '__main__':
    main()