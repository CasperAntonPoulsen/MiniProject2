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

type Connection struct { // Creating a struct to contain each connection from clients
	stream pb.ChittyChat_CreateStreamServer // Each connection contains a stream object which is used to send the message during the broadcast function
	id     string
	active bool // A bool used to check if a connection is still active
	error  chan error
}

type Server struct {
	pb.UnimplementedChittyChatServer
	Connection []*Connection
}

func (s *Server) CreateStream(pconn *pb.Connect, stream pb.ChittyChat_CreateStreamServer) error { // Once a user wants to connect, we create a stream object and add their usercredentials to it in a connection stuct
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn) // connection is added to the server

	return <-conn.error // channel any errors out
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *pb.ChatMessage) (*pb.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *pb.ChatMessage, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				if err != nil {
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &pb.Empty{}, nil
}

func main() {
	//var connections []*Connection

	//server := &Server{connections}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Error, couldn't create the server %v", err)
	}

	log.Printf("Starting server at port %v", port)

	pb.RegisterChittyChatServer(grpcServer, &Server{})
	grpcServer.Serve(listener)
}
