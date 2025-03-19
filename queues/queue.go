package queues

import (
	pb "github.com/franveiga/MalatoMQ/protos"
)

type Queue interface {
	QueueMessage(item QItem) error
    SendMessage(stream chan *pb.Message) error
	GetName() string
	GetQueue() []QItem
}
