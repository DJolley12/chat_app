package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/DJolley12/chat_app/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChatClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SendMessage(ctx, &pb.ChatMessage{
		FromName:    "Me",
		ToNames:     []string{"Me"},
		MessageBody: "a message that is cool",
		IsEncrypted: false,
	})


	log.Printf("Greeting: %s", r.GetName())
}
