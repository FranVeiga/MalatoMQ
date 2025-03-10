package queues

import (
	pb "github.com/franveiga/MalatoMQ/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type QItem struct {
	Timestamp time.Time
	Message   string
}

func NewQItem(message string) QItem {
	return QItem{
		Timestamp: time.Now(),
		Message:   message,
	}
}

func NewQItemWithTime(message string, timestamp time.Time) QItem {
	return QItem{
		Timestamp: timestamp,
		Message:   message,
	}
}

func (q QItem) ToGRPCMessage() pb.Message {
	return pb.Message{
		Content:   q.Message,
		Timestamp: timestamppb.New(q.Timestamp),
	}
}
