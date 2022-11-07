package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/DJolley12/chat_app/protos"
	"golang.org/x/sync/errgroup"
)

type client struct {
	cancel    context.CancelFunc
	chat      pb.ChatClient
	conn      *pb.Connect
	ctx       context.Context
	me        *pb.User
	inChats   chan *pb.ChatMessage
	outChats  chan *pb.ChatMessage
	recCh     chan *pb.ReceivedMessage
	errCh     chan error
	sendErrCh chan error
	g         *errgroup.Group
	stream    pb.Chat_CreateMessageStreamClient
}

func newClient(c pb.ChatClient, user *pb.User) client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	g, ctx := errgroup.WithContext(ctx)

	return client{
		cancel:    cancel,
		chat:      c,
		conn:      &pb.Connect{User: user, Active: true},
		ctx:       ctx,
		me:        user,
		inChats:   make(chan *pb.ChatMessage),
		outChats:  make(chan *pb.ChatMessage),
		recCh:     make(chan *pb.ReceivedMessage),
		errCh:     make(chan error),
		sendErrCh: make(chan error),
		g:         g,
	}
}

func (c *client) connect() error {
	stream, err := c.chat.CreateMessageStream(c.ctx, c.conn)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	c.stream = stream
	return nil
}

func (c *client) recvLoop() error {
	c.g.Go(func() error {
		for {
			msg, err := c.stream.Recv()
			if err != nil {
				return err
			}

			c.inChats <- msg
		}
	})

	go func() {
		if err := c.g.Wait(); err != nil {
			c.errCh <- err
		}
	}()

	return <-c.errCh
}

func (c *client) sendLoop() error {
	if c.stream == nil {
		return fmt.Errorf("not connected to server")
	}

	c.g.Go(func() error {
		for {
			rec, err := c.chat.SendMessage(c.ctx, <-c.outChats)
			if err != nil {
				return err
			}

			c.recCh <- rec
		}
	})

	go func() {
		if err := c.g.Wait(); err != nil {
			c.sendErrCh <- err
		}
	}()

	return <-c.sendErrCh
}

func (c *client) pullInboundMessage() *pb.ChatMessage {
	return <-c.inChats
}

func (c *client) queueSend(msg *pb.ChatMessage) error {
	if c.stream == nil {
		return fmt.Errorf("not connected to server")
	}

	c.outChats <- msg

	return nil
}
