package identityinfra

import (
	"context"
	identity "play-ground/software_acrh/master_worker/internal/worker/core/identity/entity"
)

type k8sRuntimeInfoProvider struct{}

func NewK8sRuntimeInfoProvider() *k8sRuntimeInfoProvider {
	return &k8sRuntimeInfoProvider{}
}

func (s *k8sRuntimeInfoProvider) GetMetadata(ctx context.Context) (identity.AgentMetadata, error) {
	return identity.AgentMetadata{
		Namespace:  "agent-system",
		NodeName:   "node-worker-1",
		PodName:    "agent-afs",
		Hostname:   "hostname",
		K8SVersion: "1.35.2",
	}, nil
}
