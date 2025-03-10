package queues

import (
	pb "github.com/franveiga/MalatoMQ/protos"
)

type Queue interface {
	QueueMessage(item QItem) error
	SendMessage(stream pb.MQ_ConsumeMessageServer) error
	GetName() string
	GetQueue() []QItem
}
