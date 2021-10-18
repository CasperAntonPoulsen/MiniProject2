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

func receiveMessages(ctx context.Context, client pb.ChittyChatClient) {
	cs, err := client.ReceiveMsg(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("%v.ReceiveMsg(_) = _, %v", client, err)
	}

	for {
		message, err := cs.Recv()

		if err == nil {
			log.Print(message)
			//log.Printf("%v: %v", message.GetFrom(), message.GetChatmessage())
		}
	}

}

func sendMessages(ctx context.Context, client pb.ChittyChatClient) {
	for {
		var msg string
		fmt.Scanln(&msg)
		log.Print("This message was read: " + msg)
		client.SendMsg(ctx, &pb.ChatMessage{From: name, Chatmessage: msg, Logicaltimes: logictime})
	}
}

func main() {
	ctx := context.Background()
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChittyChatClient(conn)

	fmt.Println("Enter your name: ")

	fmt.Scanln(&name)
	clientId = "5"

	j, err := c.Join(ctx, &pb.User{Id: clientId, Name: name})

	log.Print(j, err)

	users, _ := c.GetAllUsers(ctx, &pb.Empty{})

	log.Print(users)

	go receiveMessages(ctx, c)
	go sendMessages(ctx, c)

	for {

	}

}

// x := 15

//a := &x memory address 0x...

//b := *x værdien af x
