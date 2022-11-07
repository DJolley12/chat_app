package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/DJolley12/chat_app/protos"
	glog "google.golang.org/grpc/grpclog"
)

type connection struct {
	stream pb.Chat_CreateMessageStreamServer
	active bool
	error  chan error
}

type server struct {
	connections map[string]*connection
	pb.ChatServer
}

func (s *server) CreateMessageStream(in *pb.Connect, stream pb.Chat_CreateMessageStreamServer) error {
	conn := &connection{
		stream: stream,
		active: true,
		error:  make(chan error),
	}
	s.connections[in.GetUser().GetId()] = conn
	fmt.Printf("incoming connection from: %#v\n", in.GetUser())
	println("connection recv")
	return <-conn.error
}

func (s *server) SendMessage(ctx context.Context, in *pb.ChatMessage) (*pb.ReceivedMessage, error) {
	log.Printf("received message from %#v, to %#v", in.GetFromUser().GetId(), in.GetToUsers())
	wg := sync.WaitGroup{}
	done := make(chan int)

	for _, u := range in.GetToUsers() {
		wg.Add(1)

		go func(msg *pb.ChatMessage, u *pb.User) {
			defer wg.Done()
			conn := s.connections[u.GetId()]
			if conn.active {
				err := conn.stream.Send(msg)
				glog.Info("sending message to user: %#v, stream: %v", u, conn.stream)

				if err != nil {
					glog.Errorf("err sending %v - %v", u.GetId(), err)
					conn.active = false
					conn.error <- err
				}
			}
		}(in, u)

	}

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	return &pb.ReceivedMessage{
		User: in.GetFromUser(),
	}, nil
}

