package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "github.com/franveiga/MalatoMQ/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartInteractiveConsumer(host, port string) {
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
		var n int
		fmt.Print("Ingrese cantidad de mensajes a consumir: ")
		fmt.Scanf("%d\n", &n)

		stream, err := client.ConsumeMessage(context.Background(), &pb.QueueName{Name: queueName})
		if err != nil {
			log.Fatalf("Error while creating ConsumeMessage stream: %v", err)
		}
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("EOF")
					break
				} else {
					log.Fatalf("Error receiving from stream: %v", err)
				}
			}
			fmt.Println(msg.GetContent())
			break // only one message
		}
	}

}

func main() {
	host := "localhost"
	port := "3000"
	StartInteractiveConsumer(host, port)
}
