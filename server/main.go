package main

import (
	"context"

	pb "github.com/DJolley12/chat_app/protos"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.MessagerServer
}

func (s *server) SendMessage(ctx context.Context, in *pb.ChatMessage) (*pb.ReceivedMessage, error) {
	return nil, nil
}

func main() {

}
