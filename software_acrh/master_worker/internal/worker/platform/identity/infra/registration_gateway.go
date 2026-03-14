package identityinfra

import (
	"context"
	"fmt"
	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
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

func (rc *registrationClient) Register(ctx context.Context) error {
	in := pb.RegisterRequest{
		RegistrationToken: "OK1",
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
