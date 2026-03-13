package infra

import (
	"context"
	"fmt"
	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"

	"github.com/google/uuid"
)

type registrationClient struct {
	grpcClient pb.AgentServiceClient
}

func NewRegistrationClient(grpcClient pb.AgentServiceClient) *registrationClient {
	if grpcClient == nil {
		panic("infra: registration client requires grpc client")
	}

	return &registrationClient{
		grpcClient: grpcClient,
	}
}

func (rc *registrationClient) Register(ctx context.Context, agentID uuid.UUID) error {
	in := pb.RegisterRequest{
		RegistrationToken: "OK",
		AgentVersion:      "get from env",
		K8SVersion:        "v1.88.0",
	}
	resp, err := rc.grpcClient.Register(ctx, &in)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}
