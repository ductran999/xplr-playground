package handler

import (
	"context"
	"log"
	"play-ground/api_arch/grpc/gen/proto/userpb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	log.Println("GetUser called, id =", req.Id)

	return &userpb.GetUserResponse{
		Id:   req.Id,
		Name: "John Doe",
	}, nil
}
