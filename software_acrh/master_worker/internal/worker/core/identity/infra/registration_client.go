package identityinfra

import (
	"context"
	"fmt"
	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	identity "play-ground/software_acrh/master_worker/internal/worker/core/identity/entity"
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

func (rc *registrationClient) Register(ctx context.Context, agent identity.Agent) error {
	in := pb.RegisterRequest{
		RegistrationToken: agent.RegistrationToken,
		AgentVersion:      agent.Version,
		Metadata: &pb.AgentMetadata{
			Namespace:  agent.Metadata.Namespace,
			NodeName:   agent.Metadata.NodeName,
			PodName:    agent.Metadata.PodName,
			Hostname:   agent.Metadata.Hostname,
			K8SVersion: agent.Metadata.K8SVersion,
		},
	}

	resp, err := rc.grpcClient.Register(ctx, &in)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}
