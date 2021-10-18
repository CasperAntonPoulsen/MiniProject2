package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/CasperAntonPoulsen/MiniProject2/proto"
)

const port = ":8080"

type chittyChatServer struct {
	pb.UnimplementedChittyChatServer
	usersInChat []*pb.User
	observers   []pb.ChittyChat_ReceiveMsgServer //gRPC server stream
	clocks      []int32
}

func (s *chittyChatServer) Join(ctx context.Context, user *pb.User) (*pb.JoinResponse, error) {
	for _, _user := range s.usersInChat {
		if user.Name == _user.Name {
			return &pb.JoinResponse{Error: 1,
					Msg: "User already exists"},
				nil
		}
	}

	s.usersInChat = append(s.usersInChat, user)
	return &pb.JoinResponse{Error: 0,
			Msg: "Success"},
		nil
}

func (s *chittyChatServer) GetAllUsers(ctx context.Context, empty *pb.Empty) (*pb.UserList, error) {
	return &pb.UserList{Users: s.usersInChat}, nil
}

func (s *chittyChatServer) ReceiveMsg(empty *pb.Empty, stream pb.ChittyChat_ReceiveMsgServer) error { //recive msg from client
	s.observers = append(s.observers, stream)
	return nil
}

func (s *chittyChatServer) SendMsg(ctx context.Context, chatMessage *pb.ChatMessage) (*pb.Empty, error) { //send all chat msgs to clients
	log.Print(len(s.observers))
	for _, observer := range s.observers {

		err := observer.Send(chatMessage)
		log.Print(err)
		log.Print(chatMessage)
	}
	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChittyChatServer(s, &chittyChatServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
