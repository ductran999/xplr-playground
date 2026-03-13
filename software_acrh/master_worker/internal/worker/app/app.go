package app

import (
	"context"
	"log/slog"
	agentv1 "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	"play-ground/software_acrh/master_worker/internal/worker/config"
	"play-ground/software_acrh/master_worker/internal/worker/registration/infra"
	"play-ground/software_acrh/master_worker/internal/worker/registration/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerApp struct {
	cfg *config.Config

	Conn        *grpc.ClientConn
	AgentClient agentv1.AgentServiceClient

	registerClusterUC usecase.RegisterClusterUseCase
}

func Initialize(cfg *config.Config) (*WorkerApp, error) {
	conn, err := grpc.NewClient(
		cfg.ServerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	agentClient := agentv1.NewAgentServiceClient(conn)

	// Registration
	registrationClient := infra.NewRegistrationClient(agentClient)
	registerClusterUC := usecase.NewRegisterClusterUC(registrationClient)

	return &WorkerApp{
		cfg:               cfg,
		Conn:              conn,
		AgentClient:       agentClient,
		registerClusterUC: registerClusterUC,
	}, nil
}

func (wa *WorkerApp) Run(ctx context.Context) error {
	if err := wa.registerClusterUC.Execute(ctx); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (wa *WorkerApp) Close() {
	if wa.Conn == nil {
		return
	}

	if err := wa.Conn.Close(); err != nil {
		slog.Warn("close grpc connection failed", "error", err)
	}
}
