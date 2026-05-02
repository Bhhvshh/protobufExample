# gRPC Practice Examples

A collection of gRPC examples using Python and Go, ranging from a simple Hello/Greeter service to bidirectional streaming chat and a full event-driven pipeline with MongoDB and RabbitMQ.

---

## Table of Contents

1. [Example 1 – Python Hello/Greeter with CDC Pipeline (root)](#example-1--python-hellogreeter-with-cdc-pipeline)
2. [Example 2 – Python Bidirectional Chat Streaming (`grpcexample/`)](#example-2--python-bidirectional-chat-streaming)
3. [Example 3 – Go Bidirectional Chat Streaming (`grpc/`)](#example-3--go-bidirectional-chat-streaming)
4. [Example 4 – Go gRPC + MongoDB + CDC + RabbitMQ Pipeline (`lpw/`)](#example-4--go-grpc--mongodb--cdc--rabbitmq-pipeline)

---

## Example 1 – Python Hello/Greeter with CDC Pipeline

**Location:** project root  
**Language:** Python  
**Pattern:** Unary RPC → MongoDB → CDC polling → RabbitMQ

### Overview

A gRPC `Greeter` service where the server stores every greeting in MongoDB. A Change Data Capture (CDC) service polls MongoDB for new records and publishes them to a RabbitMQ queue. The client sends gRPC requests and simultaneously consumes messages from the queue.

```
Client ──(gRPC SayHello)──▶ Server ──▶ MongoDB
                                           │
                                      CDC service
                                           │
                                       RabbitMQ
                                           │
Client ◀──(consume)────────────────────────┘
```

### Project Structure

```
protos/hello.proto      # Protobuf definition (Greeter service)
server.py               # gRPC server – stores greetings in MongoDB
client.py               # gRPC client – sends requests & reads from RabbitMQ
cdc_service.py          # CDC service – polls MongoDB, publishes to RabbitMQ
hello_pb2.py            # Generated protobuf code
hello_pb2_grpc.py       # Generated gRPC code
docker-compose.yml      # Starts RabbitMQ (with management UI)
requirements.txt        # Python dependencies
```

### Prerequisites

- Python 3.8+
- MongoDB running on `localhost:27017`
- RabbitMQ (start with Docker Compose)

### Setup

```bash
# Install Python dependencies
pip install -r requirements.txt

# Start RabbitMQ
docker-compose up -d
```

### Running

Open three terminals:

```bash
# Terminal 1 – start the gRPC server
python server.py

# Terminal 2 – start the CDC service
python cdc_service.py

# Terminal 3 – start the client
python client.py
```

Enter a name in the client. The server stores the greeting in MongoDB, the CDC service picks it up and publishes it to RabbitMQ, and the client prints the message received from the queue.

---

## Example 2 – Python Bidirectional Chat Streaming

**Location:** `grpcexample/`  
**Language:** Python  
**Pattern:** Bidirectional streaming RPC

### Overview

A multi-client chat application using gRPC bidirectional streaming. Each connected client can send and receive messages in real time. The server broadcasts every incoming message to all connected clients.

### Project Structure

```
grpcexample/
├── chat.proto          # Protobuf definition (ChatService – bidirectional stream)
├── server.py           # Chat server – broadcasts messages to all clients
├── client.py           # Chat client – sends & receives messages concurrently
├── chat_pb2.py         # Generated protobuf code
├── chat_pb2_grpc.py    # Generated gRPC code
└── component.jsx       # Example React component (UI skeleton)
```

### Prerequisites

- Python 3.8+
- `grpcio` and `grpcio-tools` (see `requirements.txt` in root)

### Regenerating Protobuf Code

```bash
cd grpcexample
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. chat.proto
```

### Running

```bash
# Terminal 1 – start the server
cd grpcexample
python server.py

# Terminal 2+ – start one or more clients
cd grpcexample
python client.py
```

Enter a name when prompted and start chatting. Open multiple client terminals to see messages broadcast between them.

---

## Example 3 – Go Bidirectional Chat Streaming

**Location:** `grpc/`  
**Language:** Go  
**Pattern:** Bidirectional streaming RPC

### Overview

The same multi-client chat concept as Example 2, implemented in Go. The server tracks all active streams and broadcasts every received message to every connected client.

### Project Structure

```
grpc/
├── proto/
│   └── chat.proto          # Protobuf definition (ChatService – bidirectional stream)
├── example.com/grpcdemo/   # Generated Go protobuf/gRPC code
├── client/
│   └── main.go             # Chat client
├── main.go                 # Chat server
├── go.mod
└── go.sum
```

### Prerequisites

- Go 1.21+
- `protoc` compiler and `protoc-gen-go` / `protoc-gen-go-grpc` plugins (only needed to regenerate code)

### Regenerating Protobuf Code

```bash
cd grpc
protoc --go_out=. --go-grpc_out=. proto/chat.proto
```

### Running

```bash
# Terminal 1 – start the server
cd grpc
go run main.go

# Terminal 2+ – start one or more clients
cd grpc/client
go run main.go
```

Enter a name when prompted and start chatting.

---

## Example 4 – Go gRPC + MongoDB + CDC + RabbitMQ Pipeline

**Location:** `lpw/`  
**Language:** Go  
**Pattern:** Unary RPC → MongoDB Change Streams → RabbitMQ

### Overview

A Go implementation of a full event-driven pipeline. The gRPC server persists messages to MongoDB. A CDC service uses MongoDB Change Streams to watch for inserts and publishes them to a RabbitMQ queue. The client sends messages via gRPC and consumes processed events from the queue.

```
Client ──(gRPC SendMessage)──▶ Server ──▶ MongoDB
                                              │
                                       CDC (change stream)
                                              │
                                          RabbitMQ
                                              │
Client ◀──(consume)────────────────────────────┘
```

### Project Structure

```
lpw/
├── proto/
│   ├── message.proto       # Protobuf definition (MessageService – unary RPC)
│   ├── message.pb.go       # Generated protobuf code
│   └── message_grpc.pb.go  # Generated gRPC code
├── server/
│   └── main.go             # gRPC server – stores messages in MongoDB
├── client/
│   └── main.go             # gRPC client – sends messages & reads from RabbitMQ
├── cdc/
│   └── main.go             # CDC service – watches MongoDB, publishes to RabbitMQ
├── lpw.proto               # High-level design notes
└── go.mod
```

### Prerequisites

- Go 1.21+
- MongoDB running on `localhost:27017` (must have replica set enabled for Change Streams)
- RabbitMQ running on `localhost:5672`

### Regenerating Protobuf Code

```bash
cd lpw
protoc --go_out=. --go-grpc_out=. proto/message.proto
```

### Running

Open three terminals:

```bash
# Terminal 1 – start the gRPC server
cd lpw/server
go run main.go

# Terminal 2 – start the CDC service
cd lpw/cdc
go run main.go

# Terminal 3 – start the client
cd lpw/client
go run main.go
```

The client sends a test message via gRPC, the server writes it to MongoDB, the CDC service detects the change and publishes it to RabbitMQ, and the client prints the message received from the queue.

---

## Dependencies

### Python (root & `grpcexample/`)

| Package | Version |
|---|---|
| grpcio | 1.59.0 |
| grpcio-tools | 1.59.0 |
| protobuf | 3.20.3 |
| pymongo | 4.6.0 |
| pika | 1.3.2 |

Install with:
```bash
pip install -r requirements.txt
```

### Go (`grpc/` and `lpw/`)

Dependencies are managed with Go modules. Run `go mod tidy` inside the respective directory to fetch them.

| Package | Used in |
|---|---|
| google.golang.org/grpc | `grpc/`, `lpw/` |
| google.golang.org/protobuf | `grpc/`, `lpw/` |
| go.mongodb.org/mongo-driver | `lpw/` |
| github.com/streadway/amqp | `lpw/` |
