# Simple gRPC Python Example

This is a simple example of a gRPC client-server application in Python.

## Setup

1. Make sure you have Python installed
2. Install the required packages:
   ```bash
   pip install -r requirements.txt
   ```

## Running the Application

1. First, start the server:

   ```bash
   python server.py
   ```

2. In a new terminal, run the client:

   ```bash
   python client.py
   ```

3. When prompted, enter your name and see the server's response.

## Project Structure

- `protos/hello.proto`: The protocol buffer definition file
- `server.py`: The gRPC server implementation
- `client.py`: The gRPC client implementation
- `hello_pb2.py`: Generated protocol buffer code
- `hello_pb2_grpc.py`: Generated gRPC code
