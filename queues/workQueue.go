package queues

import (
	pb "github.com/franveiga/MalatoMQ/protos"
)

type WorkQueue struct {
	queue
	// TODO: consider adding a queue of waiters to respect FIFO on multiple consumers
}

func NewWorkQueue(name string) WorkQueue {
	return WorkQueue{queue: newQueue(name)}
}

func (wq *WorkQueue) QueueMessage(item QItem) error {
	wq.Queue(item)
	return nil
}

func (wq *WorkQueue) SendMessage(stream chan *pb.Message) error {
	msg := wq.Dequeue().ToGRPCMessage()
	stream<-&msg
    close(stream) // Message queues only send one message
	return nil
}

func (wq *WorkQueue) GetName() string {
	return wq.Name
}

func (wq *WorkQueue) GetQueue() []QItem {
	return wq.Messages
}
