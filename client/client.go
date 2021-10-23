package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "github.com/CasperAntonPoulsen/MiniProject2/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

var client pb.ChittyChatClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *pb.User) error {
	var streamerror error
	stream, err := client.CreateStream(context.Background(), &pb.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("connection has failed: %v", err)
	}
	wait.Add(1)
	go func(str pb.ChittyChat_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("error reading message: %v", err)
				break
			}
			fmt.Printf("%v %v : %s : %v\n", msg.Id, msg.From, msg.Chatmessage, msg.GetLogicaltimes())
		}
	}(stream)
	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make(chan int)

	name := flag.String("Name", "Anon", "The name of the user")
	flag.Parse()

	id := sha256.Sum256([]byte(timestamp.String() + *name))

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect to service : %v", err)
	}

	client = pb.NewChittyChatClient(conn)
	user := &pb.User{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	connect(user)

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &pb.ChatMessage{
				Id:          user.Id,
				From:        user.Name,
				Chatmessage: scanner.Text(),
			}

			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("error sending message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
