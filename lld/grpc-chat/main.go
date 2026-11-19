package main

import (
	"fmt"
	"log"
	"net"
	"system-design/lld/grpc-chat/gen"
	"system-design/lld/grpc-chat/handler"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	var conn []*handler.Connection

	pool := &handler.Pool{
		Connection: conn,
	}

	gen.RegisterBroadcastServer(grpcServer, pool)

	listner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	fmt.Println("Server started at port :8080")

	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
}
