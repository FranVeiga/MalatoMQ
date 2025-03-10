package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/franveiga/MalatoMQ/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newMessage(queueName, msg string) *pb.Message {
	ret := pb.Message{
		Content:   msg,
		Timestamp: timestamppb.New(time.Now()),
		Queue:     queueName,
	}
	return &ret
}

func StartInteractiveClient(host, port string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(fmt.Sprintf("%v:%v", host, port), opts...)
	if err != nil {
		log.Fatalf("Unable to connect to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewMQClient(conn)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Name of the queue: ")
	queueName, _ := reader.ReadString('\n')
	queueName = strings.Trim(queueName, "\n")

	for {
		fmt.Print("Msg to send: ")
		s, _ := reader.ReadString('\n')
		s = strings.TrimSpace(s)
		fmt.Printf("Valor de s: %v\n", s)
		if s == "q\n" {
			break
		}
		stream, err := client.SendMessage(context.Background())
		if err != nil {
			log.Fatalf("Error while creating SendMessage stream: %v", err)
		}
		for subs := range strings.SplitSeq(s, ";") {
			stream.Send(newMessage(queueName, subs))
		}
		err = stream.CloseSend()
		if err != nil {
			log.Fatalf("Error closing send stream: %v", err)
		}
	}

}

func main() {
	host := "localhost"
	port := "3000"
	StartInteractiveClient(host, port)
}
