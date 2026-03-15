package monitoringinfra

import (
	"context"
	agentv1 "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	"play-ground/software_acrh/master_worker/internal/worker/core/monitoring/domain"
	monitoringuc "play-ground/software_acrh/master_worker/internal/worker/core/monitoring/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type statusReporter struct {
	grpcClient agentv1.AgentServiceClient
}

func NewStatusReporter(grpcClient agentv1.AgentServiceClient) monitoringuc.StatusReporter {
	if grpcClient == nil {
		panic("nil grpc client")
	}

	return &statusReporter{
		grpcClient: grpcClient,
	}
}

func (s *statusReporter) SendHeartbeat(ctx context.Context, agentID string, state string) error {
	req := &agentv1.SendHeartbeatRequest{
		Status:  state,
		AgentId: agentID,
	}

	_, err := s.grpcClient.SendHeartbeat(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		if st.Code() == codes.Unauthenticated {
			return domain.ErrAgentUnauthorized
		}

		return err
	}

	return nil
}
