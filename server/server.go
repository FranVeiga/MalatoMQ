package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/franveiga/MalatoMQ/protos"
	"github.com/franveiga/MalatoMQ/queues"
	"google.golang.org/grpc"
)

var (
	host = flag.String("h", "localhost", "The network host the server will listen on")
	port = flag.String("p", "3000", "The port the server will listen on")
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", *host, *port))
	if err != nil {
		log.Fatalf("Couldn't create listener on %v:%v\n%v", host, port, err)
	}
	defer lis.Close()

	server := grpc.NewServer()

	queuesList := make([]queues.Queue, 0)
	wq := queues.NewWorkQueue("work_queue1")
	queuesList = append(queuesList, &wq)

	msgQueueServer := newMQServer(queuesList)
	pb.RegisterMQServer(server, &msgQueueServer)
	fmt.Printf("Listening on %v:%v\n", *host, *port)
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Stop()
}
