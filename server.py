import grpc
from concurrent import futures
import hello_pb2
import hello_pb2_grpc
from pymongo import MongoClient
from datetime import datetime

class Greeter(hello_pb2_grpc.GreeterServicer):
    def __init__(self):
        # Connect to MongoDB
        self.client = MongoClient('mongodb://localhost:27017/')
        self.db = self.client['grpc_messages']
        self.collection = self.db['greetings']

    def SayHello(self, request, context):
        # Create message document
        message_doc = {
            'name': request.name,
            'message': f"Hello, {request.name}!",
            'timestamp': datetime.now()
        }
        
        # Save to MongoDB
        self.collection.insert_one(message_doc)
        
        return hello_pb2.HelloResponse(message=f"Hello, {request.name}!")

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    hello_pb2_grpc.add_GreeterServicer_to_server(Greeter(), server)
    server.add_insecure_port('[::]:50051')
    print("Server starting on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()