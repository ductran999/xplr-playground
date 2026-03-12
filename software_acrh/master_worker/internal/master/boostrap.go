package master

import (
	"fmt"
	"log"
	"net"
	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"

	"google.golang.org/grpc"
)

func Run() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	connManager := NewConnectionManager()
	hdl := NewHandler(connManager)
	pb.RegisterAgentServiceServer(grpcServer, hdl)

	log.Println("Control Plane listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
