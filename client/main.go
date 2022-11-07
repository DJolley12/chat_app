package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	pb "github.com/DJolley12/chat_app/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	glog "google.golang.org/grpc/grpclog"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type serverStatus int

const (
	connected serverStatus = iota + 1
	// reconnecting
	fatal
)

func main() {
	// flag.Parse()
	println("post flag")
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewChatClient(conn)
	client := newClient(c, &pb.User{Id: "1", Name: "Me"})
	defer client.cancel()

	println("before conn loop")
	sCh := make(chan serverStatus, 0)
	defer close(sCh)

	go connectionLoop(&client, sCh)

	go sendLoop(&client, sCh)

	chain := newMessageChain()

	go inboundLoop(&client, &chain)

	for {
		userInputLoop(&client)
	}

	// Contact the server and print out its response.
	// r, err := client.SendMessage(ctx, &pb.ChatMessage{
	// 	FromName:    "Me",
	// 	ToNames:     []string{"Me"},
	// 	MessageBody: "a message that is cool",
	// 	IsEncrypted: false,
	// })
}

func userInputLoop(client *client) {
	var input string
	println("Message:")
	fmt.Scanln(&input)
	msg := &pb.ChatMessage{
		FromUser: client.me,
		ToUsers: []*pb.User{client.me},
		MessageBody: input,
		IsEncrypted: false,
	}
	client.queueSend(msg)
}

func connectionLoop(client *client, sCh chan serverStatus) {
	println("conn loop")
	for {
		var input string
		if err := client.connect(); err != nil {
			glog.Error("cannot connect to server:", err)
			fmt.Println("do you want to reconnect?")
			fmt.Println("Y/N")
			fmt.Scanln(&input)

			if input == strings.ToLower("N") {
				sCh <- fatal
				return
			} else if input == strings.ToLower("Y") {
				continue
			} else {
				fmt.Println("did not recognize input:", &input)
			}
		}

		sCh <- connected

		if err := client.recvLoop(); err != nil {
			glog.Error("error receiving message:", err)
		}
	}
}

func inboundLoop(client *client, chain *messageChain) {
	chain.add(client.pullInboundMessage())
	for _, m := range chain.get(0, 20) {
		glog.Printf("%v\n", m)
	}
}

func sendLoop(client *client, sCh chan serverStatus) {
	status := <-sCh
	if status == connected {
		if err := client.sendLoop(); err != nil {
			glog.Error("error sending message:", err)
		}
	} else if status == fatal {
		return
	}
}
