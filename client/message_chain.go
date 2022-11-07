package main

import (
	"sync"

	pb "github.com/DJolley12/chat_app/protos"
)

type messageChain struct {
	chain  []*pb.ChatMessage
	mu sync.RWMutex
}

func newMessageChain() messageChain {
	return messageChain{
		chain: make([]*pb.ChatMessage, 0),
	}
}

func (c *messageChain) add(msg *pb.ChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.chain = append(c.chain, msg)
}

func (c *messageChain) get(ind, length int) []*pb.ChatMessage {
	c.mu.RLock()
	defer c.mu.RLock()

	low := ind
	if low > len(c.chain) {
		low = len(c.chain)
	}

	hi := length
	if length > len(c.chain) {
		hi = len(c.chain)
	}

	return c.chain[low:hi]
}
