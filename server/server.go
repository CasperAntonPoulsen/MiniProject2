package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/CasperAntonPoulsen/MiniProject2/proto"
	"google.golang.org/grpc"
)

const port = ":8080"

var mutex sync.Mutex

type Connection struct { // Creating a struct to contain each connection from clients
	stream      pb.ChittyChat_CreateStreamServer // Each connection contains a stream object which is used to send the message during the broadcast function
	id          string
	name        string
	active      bool // A bool used to check if a connection is still active
	logicaltime int32
	error       chan error
}

type Server struct {
	pb.UnimplementedChittyChatServer
	Connection  []*Connection
	logicaltime int32
}

func (s *Server) CreateStream(pconn *pb.Connect, stream pb.ChittyChat_CreateStreamServer) error { // Once a user wants to connect, we create a stream object and add their usercredentials to it in a connection stuct
	conn := &Connection{
		stream:      stream,
		id:          pconn.User.Id,
		name:        pconn.User.Name,
		active:      true,
		logicaltime: 1, // Logical time for a new connection starts at 1, because it's inffered they would have to send a connection message to call this function
		error:       make(chan error),
	}

	s.Connection = append(s.Connection, conn) // connection is added to the server

	_, err := s.BroadcastMessage(context.Background(), &pb.ChatMessage{
		Id:          "server",
		Chatmessage: pconn.User.Name + " has connected",
	})
	if err != nil {
		log.Printf("could send connection message: %v", err)
	}

	return <-conn.error // channel any errors out
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *pb.ChatMessage) (*pb.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	times := s.CalculateTime(msg)
	ClientIndex := 1
	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *pb.ChatMessage, conn *Connection) {
			defer wait.Done()

			if conn.active {
				s.logicaltime += 1

				mutex.Lock() // lock

				times[0] += 1

				msg.Logicaltimes = times

				msg.ClientTime = int32(ClientIndex)
				err := conn.stream.Send(msg)
				if err != nil {
					mutex.Unlock() //unlock
					conn.active = false
					conn.error <- err

					_, err := s.BroadcastMessage(context.Background(), &pb.ChatMessage{
						Id:          "server",
						Chatmessage: conn.name + " has disconnected",
					})
					if err != nil {
						log.Printf("could send disconnect message: %v", err)
					}
					return

				}

				mutex.Unlock() //unlock
			}
		}(msg, conn)
		ClientIndex++
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &pb.Empty{}, nil
}

func (s *Server) CalculateTime(msg *pb.ChatMessage) []int32 {
	var times []int32

	times = append(times, s.logicaltime)

	for _, conn := range s.Connection {
		if msg.GetId() == conn.id {
			conn.logicaltime += 1
		}
		times = append(times, conn.logicaltime)
	}

	return times
}

func main() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error, couldn't create the server %v", err)
	}

	log.Printf("Starting server at port %v", port)

	pb.RegisterChittyChatServer(grpcServer, &Server{})
	grpcServer.Serve(listener)
}
