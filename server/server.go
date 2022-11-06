package main

import (
	"context"
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
	return <-conn.error
}

func (s *server) SendMessage(ctx context.Context, in *pb.ChatMessage) (*pb.ReceivedMessage, error) {
	log.Printf("received message from %v, to %v", in.GetFromName(), in.GetToNames())
	wg := sync.WaitGroup{}
	done := make(chan int)

	for _, id := range in.GetToNames() {
		wg.Add(1)

		go func(msg *pb.ChatMessage, id string) {
			defer wg.Done()
			conn := s.connections[id]
			if conn.active {
				err := conn.stream.Send(msg)
				glog.Info("sending message to id: %v, stream: %v", id, conn.stream)

				if err != nil {
					glog.Errorf("err sending %v - %v", id, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(in, id)

	}

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	return &pb.ReceivedMessage{Name: in.GetFromName()}, nil
}

