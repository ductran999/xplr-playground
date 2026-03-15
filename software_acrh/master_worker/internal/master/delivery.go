package master

import (
	"context"
	"log"
	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	"strings"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedAgentServiceServer

	connectionManager *ConnectionManager
}

func NewHandler(connManger *ConnectionManager) *grpcHandler {
	return &grpcHandler{
		connectionManager: connManger,
	}
}

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

func (s *grpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Agent registering: %v", req)

	if req.RegistrationToken != "OK" {
		return nil, status.Error(codes.Unauthenticated, "invalid registration token")
	}

	return &pb.RegisterResponse{
		ClusterId:     "cluster-uuid-123",
		AgentIdentity: "secure-jwt-or-cert-id",
		Message:       "Registration successful",
	}, nil
}

func (s *grpcHandler) SendHeartbeat(ctx context.Context, req *pb.SendHeartbeatRequest) (*pb.SendHeartbeatResponse, error) {
	log.Println(req)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid token format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	log.Println("token:", token)

	return &pb.SendHeartbeatResponse{
		NextIntervalSeconds: 10,
	}, nil
}

func (s *grpcHandler) ConnectTunnel(stream pb.AgentService_ConnectTunnelServer) error {
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
