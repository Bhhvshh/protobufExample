import grpc
import hello_pb2
import hello_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = hello_pb2_grpc.GreeterStub(channel)
        name = input("Enter your name: ")
        response = stub.SayHello(hello_pb2.HelloRequest(name=name))
        print("Server response:", response.message)

if __name__ == '__main__':
    run()