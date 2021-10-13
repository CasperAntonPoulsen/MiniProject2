package main

import (
	"context"

	"github.com/golang/protobuf/proto"

	pb "github.com/CasperAntonPoulsen/MiniProject2/proto"
)

const port = ":8080"

type chittyChatServer struct {
	pb.UnimplementedChittyChatServer
	usersInChat []*pb.User
	observers []
}


var clocks []int32


func (s *chittyChatServer) Join(ctx context.Context, user *pb.User) (*pb.JoinResponse, error) {
	for _, _user := range s.usersInChat {
		if proto.Equal(user.Name, _user.Name) {
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
	return &pb.UserList{Users:s.usersInChat}, nil
}

func (s *chittyChatServer) receiveMsg(empty *pb.Empty, stream pb.ChittyChat_ReceiveMsgServer) error {

	return nil
}

func (s *chittyChatServer) sendMsg(ctx context.Context, )