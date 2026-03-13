package usecase

import (
	"context"

	"github.com/google/uuid"
)

type RegisterClusterUseCase interface {
	Execute(ctx context.Context) error
}

type RegistrationClient interface {
	Register(ctx context.Context, agentID uuid.UUID) error
}

type registerClusterUC struct {
	registrationClient RegistrationClient
}

func NewRegisterClusterUC(registrationClient RegistrationClient) RegisterClusterUseCase {
	return &registerClusterUC{
		registrationClient: registrationClient,
	}
}

func (uc *registerClusterUC) Execute(ctx context.Context) error {
	agentID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	return uc.registrationClient.Register(ctx, agentID)
}
