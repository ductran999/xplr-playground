package main

import (
	"log"
	"net"
	"play-ground/api_arch/grpc/gen/proto/streampb"
	"play-ground/api_arch/grpc/gen/proto/userpb"
	"play-ground/api_arch/grpc/handler"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &handler.UserServer{})
	streampb.RegisterProvisionServiceServer(grpcServer, &handler.WatchProvision{})

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
