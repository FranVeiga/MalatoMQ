package queues

import (
	pb "github.com/franveiga/MalatoMQ/protos"
)

type WorkQueue struct {
	genericQueue
	// TODO: consider adding a queue of waiters to respect FIFO on multiple consumers
}

func NewWorkQueue(name string) WorkQueue {
	return WorkQueue{genericQueue: newGenericQueue(name)}
}

func (wq *WorkQueue) QueueMessage(item QItem) error {
	wq.Queue(item)
	return nil
}

func (wq *WorkQueue) SendMessage(stream pb.MQ_ConsumeMessageServer) error {
	msg := wq.Dequeue().ToGRPCMessage()
	err := stream.Send(&msg)
	return err
}

func (wq *WorkQueue) GetName() string {
	return wq.Name
}

func (wq *WorkQueue) GetQueue() []QItem {
	return wq.Messages
}
