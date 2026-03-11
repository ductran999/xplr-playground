package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"

	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ConnectionManager struct {
	mu      sync.RWMutex
	tunnels map[string]chan *pb.ConnectTunnelResponse
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		tunnels: make(map[string]chan *pb.ConnectTunnelResponse),
	}
}

func (cm *ConnectionManager) Register(clusterID string, cmdChan chan *pb.ConnectTunnelResponse) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.tunnels[clusterID] = cmdChan
}

func (cm *ConnectionManager) Unregister(clusterID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if ch, ok := cm.tunnels[clusterID]; ok {
		close(ch)
	}
	delete(cm.tunnels, clusterID)
}

func (cm *ConnectionManager) SendCommand(clusterID string, cmd *pb.ConnectTunnelResponse) bool {
	cm.mu.RLock()
	ch, ok := cm.tunnels[clusterID]
	cm.mu.RUnlock()

	if !ok {
		return false
	}

	select {
	case ch <- cmd:
		return true
	default:
		return false
	}
}

type server struct {
	pb.UnimplementedAgentServiceServer

	connectionManager *ConnectionManager
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Agent registering: %s (Version: %s)", req.Hostname, req.AgentVersion)

	if req.RegistrationToken != "OK" {
		return nil, errors.New("unauthorized")
	}

	return &pb.RegisterResponse{
		ClusterId:     "cluster-uuid-123",
		AgentIdentity: "secure-jwt-or-cert-id",
		Message:       "Registration successful",
	}, nil
}

func (s *server) ConnectTunnel(stream pb.AgentService_ConnectTunnelServer) error {
	clusterID, _ := getClusterIDFromMetadata(stream.Context())

	cmdChan := make(chan *pb.ConnectTunnelResponse, 100)
	testCmd := &pb.ConnectTunnelResponse{
		CommandId: "cmd-uuid-123",
		Action:    "DEPLOY_MODEL",
	}
	cmdChan <- testCmd
	s.connectionManager.Register(clusterID, cmdChan)
	defer s.connectionManager.Unregister(clusterID)

	go func() {
		for cmd := range cmdChan {
			if err := stream.Send(cmd); err != nil {
				return
			}
		}
	}()

	for {
		report, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Received report from %s: %s", clusterID, report.Status)
	}
}

func getClusterIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("x-cluster-id")
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "cluster-id not found")
	}

	return values[0], nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterAgentServiceServer(s, &server{
		connectionManager: NewConnectionManager(),
	})

	log.Println("Control Plane listening on :50051")
	s.Serve(lis)
}
