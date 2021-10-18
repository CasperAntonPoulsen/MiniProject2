package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/CasperAntonPoulsen/MiniProject2/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

var (
	clientId  string
	name      string
	logictime []int32
)

func receiveMessages(client pb.ChittyChatClient) {
	cs, err := client.ReceiveMsg(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("%v.ReceiveMsg(_) = _, %v", client, err)
	}

	for {
		message, err := cs.Recv()
		if err == nil {
			log.Printf("%v: %v", message.GetFrom(), message.GetChatmessage())
		}
	}

}

func sendMessages(client pb.ChittyChatClient) {
	for {
		var msg string
		fmt.Scanln(&msg)

		client.SendMsg(context.Background(), &pb.ChatMessage{From: name, Chatmessage: msg, Logicaltimes: logictime})
	}
}

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChittyChatClient(conn)

	fmt.Println("Enter your name: ")

	fmt.Scanln(&name)
	clientId = "5"

	j, err := c.Join(context.Background(), &pb.User{Id: clientId, Name: name})

	log.Print(j, err)

	go receiveMessages(c)
	go sendMessages(c)

}
