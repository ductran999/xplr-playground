package identityuc

import (
	"context"
)

type RegisterClusterUseCase interface {
	Execute(ctx context.Context) error
}

type RegistrationGateway interface {
	Register(ctx context.Context) error
}

type registerClusterUC struct {
	gateway RegistrationGateway
}

func NewRegisterClusterUC(gateway RegistrationGateway) RegisterClusterUseCase {
	return &registerClusterUC{
		gateway: gateway,
	}
}

func (uc *registerClusterUC) Execute(ctx context.Context) error {
	return uc.gateway.Register(ctx)
}
