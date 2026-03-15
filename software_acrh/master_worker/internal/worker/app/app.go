package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	agentv1 "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	identityinfra "play-ground/software_acrh/master_worker/internal/worker/core/identity/infra"
	identityuc "play-ground/software_acrh/master_worker/internal/worker/core/identity/usecase"
	"play-ground/software_acrh/master_worker/internal/worker/core/monitoring/domain"
	monitoringinfra "play-ground/software_acrh/master_worker/internal/worker/core/monitoring/infra"
	monitoringuc "play-ground/software_acrh/master_worker/internal/worker/core/monitoring/usecase"
	"play-ground/software_acrh/master_worker/internal/worker/platform/config"
	gclient "play-ground/software_acrh/master_worker/internal/worker/platform/grpc"
)

type WorkerApp struct {
	cfg  *config.Config
	Conn *grpc.ClientConn

	AgentClient       agentv1.AgentServiceClient
	registerClusterUC identityuc.RegisterClusterUseCase
	monitoringUC      monitoringuc.MonitoringUseCase
}

func Initialize(cfg *config.Config) (*WorkerApp, error) {
	conn, err := gclient.Connect(gclient.Config{
		Address: cfg.ServerURL,
	})
	if err != nil {
		return nil, err
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	agentClient := agentv1.NewAgentServiceClient(conn)

	// Registration
	runtimeInfoProvider := identityinfra.NewK8sRuntimeInfoProvider()
	registrationClient := identityinfra.NewRegistrationClient(agentClient)
	registerClusterUC := identityuc.NewRegisterClusterUC(cfg, registrationClient, runtimeInfoProvider)

	// Monitoring
	statusReporter := monitoringinfra.NewStatusReporter(agentClient)
	monitoringUC := monitoringuc.NewMonitoringUseCase(statusReporter)

	return &WorkerApp{
		cfg:               cfg,
		Conn:              conn,
		AgentClient:       agentClient,
		registerClusterUC: registerClusterUC,
		monitoringUC:      monitoringUC,
	}, nil
}

func (wa *WorkerApp) Run(ctx context.Context) error {
	if err := wa.registerClusterUC.Execute(ctx); err != nil {
		return err
	}
	slog.Info("registration completed successfully!")

	go wa.runHeartbeat(ctx, "agent-1", "cluster-123")

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

func (wa *WorkerApp) runHeartbeat(ctx context.Context, agentID string, clusterID string) {
	for {
		log.Printf("Cluster %s: starting heartbeat loop...", clusterID)

		md := metadata.Pairs("x-cluster-id", clusterID)
		ctx := metadata.NewOutgoingContext(ctx, md)

		ticker := time.NewTicker(10 * time.Second)
		consecutiveSuccess := 0

		for {
			select {
			case <-ctx.Done():
				log.Printf("Cluster %s: heartbeat loop exiting due to context cancellation", clusterID)
				ticker.Stop()
				return

			case <-ticker.C:
				err := wa.monitoringUC.SendHeartbeat(ctx, agentID)
				if err != nil {
					if errors.Is(err, domain.ErrAgentUnauthorized) {
						slog.Error("Critical authentication failure",
							"agent_id", agentID,
							"cluster_id", clusterID,
							"err_type", "unauthorized",
							"action", "heartbeat_stopped",
						)
						return
					}
					slog.Error("Heartbeat failed", "error", err)
				}

				consecutiveSuccess++
				slog.Debug("Heartbeat OK", "count", consecutiveSuccess)

				if consecutiveSuccess == 1 || consecutiveSuccess%100 == 0 {
					slog.Info("Heartbeat OK", "count", consecutiveSuccess)
				}
			}
		}
	}
}
