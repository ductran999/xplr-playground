package identityuc

import (
	"context"
	identity "play-ground/software_acrh/master_worker/internal/worker/core/identity/entity"
)

type RegistrationGateway interface {
	Register(ctx context.Context, agent identity.Agent) error
}

type RuntimeInfoProvider interface {
	GetMetadata(ctx context.Context) (identity.AgentMetadata, error)
}
