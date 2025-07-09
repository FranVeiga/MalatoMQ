package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	pb "github.com/franveiga/MalatoMQ/protos"
	"github.com/franveiga/MalatoMQ/queues"
)

type QItem = queues.QItem

type MQServer struct {
	pb.MQServer
	msgQueue  []QItem
	mu        sync.Mutex
	available chan bool
	queues    []queues.ServerQueue
}

func newMQServer(queuesList []queues.ServerQueue) MQServer {
	return MQServer{
		msgQueue:  make([]QItem, 1),
		available: make(chan bool, 100000), // buffer up to 100000
		queues:    queuesList,
	}
}

func (s *MQServer) SendMessage(stream pb.MQ_SendMessageServer) error {
	var err error
	msg, err := stream.Recv()
	for err == nil {
		queue, findQueueErr := s.findQueue(msg.GetQueue())
		if findQueueErr != nil {
			return findQueueErr
		}

		if msg.GetContent() == "queue" {
			fmt.Println(queue.GetQueue())
		} else {
			item := queues.NewQItemWithTime(msg.GetContent(), msg.GetTimestamp().AsTime())
			queue.QueueMessage(item)
			fmt.Println("Queued:", item)
		}
		msg, err = stream.Recv()
	}
	if err == io.EOF {
		stream.SendAndClose(&pb.Response{Ok: true})
		return nil
	} else {
		log.Println("ERROR: Unexpected end of stream,", err)
		stream.SendAndClose(&pb.Response{Ok: false})
		return err
	}
}

func (s *MQServer) ConsumeMessage(queueName *pb.QueueName, stream pb.MQ_ConsumeMessageServer) error {
	queue, err := s.findQueue(queueName.GetName())
	if err != nil {
		return err
	}
    
    msg_chan := make(chan *pb.Message)

	go queue.SendMessage(msg_chan)

    for {
        msg, more := <- msg_chan
        if more {
            stream.Send(msg)
        } else {
            break
        }
    }

	if err != nil {
		log.Fatalf("Failed sending message to consumer: %v", err)
	}
	return nil
}

// Pop the first element from the queue. Blocks if no element is present until there is at least one.
// func (s *MQServer) nextMsg() QItem {
// 	// Wait until there is something in the queue
// 	<-s.available
//
// 	if len(s.msgQueue) < 0 {
// 		log.Fatal("WTF")
// 	}
//
// 	item := s.msgQueue[0]
// 	s.msgQueue = s.msgQueue[1:]
// 	return item
// }

func (s *MQServer) findQueue(name string) (queues.ServerQueue, error) {
	for _, q := range s.queues {
		if q.GetName() == name {
			return q, nil
		}
	}
	return nil, errors.New("Queue not found")
}
