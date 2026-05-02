import grpc
import time
from concurrent import futures
import chat_pb2
import chat_pb2_grpc

# Store all active client queues
clients = []

class ChatService(chat_pb2_grpc.ChatServiceServicer):
    def ChatStream(self, request_iterator, context):
        global clients
        client_queue = []
        clients.append(client_queue)

        try:
            # Run in background: receive messages from this client and broadcast to others
            def handle_incoming():
                for new_msg in request_iterator:
                    print(f"[{new_msg.user}] {new_msg.message}")
                    # Broadcast to all clients (including self)
                    for q in clients:
                        q.append(new_msg)

            import threading
            threading.Thread(target=handle_incoming, daemon=True).start()

            # Continuously yield messages to this client
            while True:
                if client_queue:
                    yield client_queue.pop(0)
                time.sleep(0.1)

        except Exception as e:
            print("Client disconnected:", e)

        finally:
            if client_queue in clients:
                clients.remove(client_queue)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    chat_pb2_grpc.add_ChatServiceServicer_to_server(ChatService(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print(" Chat server started on port 50051")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
