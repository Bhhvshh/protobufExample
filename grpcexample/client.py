import grpc
import threading
import chat_pb2
import chat_pb2_grpc


def listen_for_messages(stub, username):
    def generator():
        while True:
            msg = input("")
            yield chat_pb2.ChatMessage(user=username, message=msg)

    responses = stub.ChatStream(generator())
    try:
        for response in responses:
            if response.user != username:  # don’t print own messages again
                print(f"\n{response.user}: {response.message}")
    except grpc.RpcError as e:
        print("Disconnected from server:", e)


def run(username):
    channel = grpc.insecure_channel("localhost:50051")
    stub = chat_pb2_grpc.ChatServiceStub(channel)

    threading.Thread(target=listen_for_messages, args=(stub, username), daemon=True).start()

    while True:
        pass  # keep alive


if __name__ == "__main__":
    name = input("Enter your name: ")
    run(name)
