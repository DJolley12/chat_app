package main

import (
	"context"
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
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.ChatServer
}

func (s *server) SendMessage(ctx context.Context, in *pb.ChatMessage) (*pb.ReceivedMessage, error) {
	log.Printf("received message from %v", in.GetFromName())
	return &pb.ReceivedMessage{Name: in.GetFromName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
