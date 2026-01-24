package main

import "fmt"

func main() {
	// ========================================
	// gRPC IN GO - High-performance RPC framework
	// ========================================

	// This is a conceptual overview. To actually run gRPC:
	// 1. Install protoc (protobuf compiler)
	// 2. Install Go plugins: protoc-gen-go, protoc-gen-go-grpc
	// 3. Define .proto files
	// 4. Generate Go code
	// 5. Implement server & client

	fmt.Println("=== gRPC Conceptual Guide ===")
	fmt.Println()

	// ========================================
	// SETUP (run these commands)
	// ========================================

	fmt.Println("1. SETUP")
	fmt.Println("   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	fmt.Println("   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
	fmt.Println()

	// ========================================
	// PROTO FILE EXAMPLE
	// ========================================

	protoExample := `
2. PROTO FILE (proto/user.proto)
--------------------------------
syntax = "proto3";

package user;
option go_package = "github.com/yourname/project/proto/user";

// The service definition
service UserService {
    // Unary RPC - single request, single response
    rpc GetUser(GetUserRequest) returns (UserResponse);

    // Server streaming - single request, stream of responses
    rpc ListUsers(ListUsersRequest) returns (stream UserResponse);

    // Client streaming - stream of requests, single response
    rpc CreateUsers(stream CreateUserRequest) returns (CreateUsersResponse);

    // Bidirectional streaming - stream both ways
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

message GetUserRequest {
    int64 id = 1;
}

message UserResponse {
    int64 id = 1;
    string name = 2;
    string email = 3;
}

message ListUsersRequest {
    int32 page_size = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message CreateUsersResponse {
    int32 created_count = 1;
}

message ChatMessage {
    string user = 1;
    string text = 2;
}
`
	fmt.Println(protoExample)

	// ========================================
	// GENERATE CODE
	// ========================================

	fmt.Println("3. GENERATE GO CODE")
	fmt.Println("   protoc --go_out=. --go-grpc_out=. proto/user.proto")
	fmt.Println()

	// ========================================
	// SERVER IMPLEMENTATION
	// ========================================

	serverExample := `
4. SERVER IMPLEMENTATION
------------------------
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "yourproject/proto/user"
)

type server struct {
    pb.UnimplementedUserServiceServer
}

// Unary RPC implementation
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    // Check context for cancellation
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }

    // In real app: fetch from database
    return &pb.UserResponse{
        Id:    req.Id,
        Name:  "Primi",
        Email: "primi@example.com",
    }, nil
}

// Server streaming implementation
func (s *server) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    users := []pb.UserResponse{
        {Id: 1, Name: "User 1", Email: "user1@test.com"},
        {Id: 2, Name: "User 2", Email: "user2@test.com"},
        {Id: 3, Name: "User 3", Email: "user3@test.com"},
    }

    for _, user := range users {
        if err := stream.Send(&user); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, &server{})

    log.Println("gRPC server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
`
	fmt.Println(serverExample)

	// ========================================
	// CLIENT IMPLEMENTATION
	// ========================================

	clientExample := `
5. CLIENT IMPLEMENTATION
------------------------
package main

import (
    "context"
    "io"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "yourproject/proto/user"
)

func main() {
    // Connect to server
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)

    // Unary call with timeout
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
    if err != nil {
        log.Fatalf("GetUser failed: %v", err)
    }
    log.Printf("User: %v", user)

    // Server streaming
    stream, err := client.ListUsers(context.Background(), &pb.ListUsersRequest{PageSize: 10})
    if err != nil {
        log.Fatalf("ListUsers failed: %v", err)
    }

    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("stream error: %v", err)
        }
        log.Printf("Streamed user: %v", user)
    }
}
`
	fmt.Println(clientExample)

	// ========================================
	// INTERCEPTORS (MIDDLEWARE)
	// ========================================

	interceptorExample := `
6. INTERCEPTORS (MIDDLEWARE)
----------------------------
// Unary interceptor (logging example)
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()

    resp, err := handler(ctx, req)

    log.Printf("Method: %s, Duration: %v, Error: %v",
        info.FullMethod, time.Since(start), err)

    return resp, err
}

// Use it
s := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor),
    // grpc.ChainUnaryInterceptor(interceptor1, interceptor2), // multiple
)
`
	fmt.Println(interceptorExample)

	// ========================================
	// PROJECT STRUCTURE
	// ========================================

	structureExample := `
7. RECOMMENDED PROJECT STRUCTURE
--------------------------------
myproject/
├── go.mod
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── proto/
│   └── user/
│       ├── user.proto
│       ├── user.pb.go        (generated)
│       └── user_grpc.pb.go   (generated)
├── internal/
│   ├── server/
│   │   └── user_service.go   (implementation)
│   └── repository/
│       └── user_repo.go
└── Makefile

# Makefile
proto:
    protoc --go_out=. --go-grpc_out=. proto/**/*.proto

run-server:
    go run cmd/server/main.go

run-client:
    go run cmd/client/main.go
`
	fmt.Println(structureExample)

	// ========================================
	// PERFORMANCE TIPS
	// ========================================

	tips := `
8. PERFORMANCE TIPS FOR 20k req/s
---------------------------------
1. Use streaming for bulk data
2. Reuse connections (gRPC does this automatically)
3. Use connection pooling on client side
4. Enable keepalive:
   grpc.KeepaliveParams(keepalive.ServerParameters{
       MaxConnectionIdle: 5 * time.Minute,
       Time:              2 * time.Hour,
   })

5. Set appropriate message size limits:
   grpc.MaxRecvMsgSize(10 * 1024 * 1024) // 10MB

6. Use context timeouts ALWAYS

7. Profile with pprof to find bottlenecks

8. Consider protobuf arena allocation for high-throughput

9. Load test with ghz:
   ghz --insecure --proto ./proto/user.proto \
       --call user.UserService.GetUser \
       -d '{"id": 1}' \
       -n 20000 -c 100 \
       localhost:50051
`
	fmt.Println(tips)

	fmt.Println("=== You're ready to build high-performance gRPC services! ===")
}
