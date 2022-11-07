package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/DJolley12/chat_app/protos"
	"google.golang.org/grpc"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("ip-address", "127.0.0.1", "ip address to serve on-defaults to localhost")
	port = flag.Int("port", 50051, "the server port-defaults to 50051")
)
func main() {
	connections := make(map[string]*connection)
	server := &server{
		connections: connections,
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServer(s, server)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
