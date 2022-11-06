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

type client struct {
	ctx context.Context
	pb.ChatClient
}

func (c *client) connect(user *pb.User) error {
	var streamerror error

	usr := &pb.User{Id: "1", Name: "Me"}
	pbConn := &pb.Connect{User: usr, Active: true}

	stream, err := client.CreateMessageStream(c.ctx, pbConn)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	stream.Recv()
}
