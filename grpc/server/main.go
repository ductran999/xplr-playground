package main

import (
	"context"
	"log"
	"net"
	"play-ground/grpc/proto/example.com/userpb"

	"google.golang.org/grpc"
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

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(
		grpcServer,
		&UserServer{},
	)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
