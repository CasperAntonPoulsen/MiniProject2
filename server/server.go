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
	observers   []pb.ChittyChat_ReceiveMsgServer
}

var clocks []int32

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

func (s *chittyChatServer) getAllUsers(ctx context.Context, empty *pb.Empty) (*pb.UserList, error) {
	return &pb.UserList{Users: s.usersInChat}, nil
}

func (s *chittyChatServer) receiveMsg(empty *pb.Empty, stream pb.ChittyChat_ReceiveMsgServer) error {
	s.observers = append(s.observers, stream)
	return nil
}

func (s *chittyChatServer) sendMsg(ctx context.Context, chatMessage *pb.ChatMessage) (*pb.Empty, error) {
	for _, observer := range s.observers {
		observer.Send(chatMessage)
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
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
