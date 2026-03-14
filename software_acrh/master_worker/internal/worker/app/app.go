package app

import (
	"context"
	"log"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	agentv1 "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	"play-ground/software_acrh/master_worker/internal/worker/platform/config"
	identityinfra "play-ground/software_acrh/master_worker/internal/worker/platform/identity/infra"
	identityuc "play-ground/software_acrh/master_worker/internal/worker/platform/identity/usecase"
	gclient "play-ground/software_acrh/master_worker/internal/worker/platform/shared/grpc"
)

type WorkerApp struct {
	cfg *config.Config

	Conn        *grpc.ClientConn
	AgentClient agentv1.AgentServiceClient

	registerClusterUC identityuc.RegisterClusterUseCase
}

func Initialize(cfg *config.Config) (*WorkerApp, error) {
	conn, err := gclient.Connect(gclient.Config{
		Address: cfg.ServerURL,
	})
	if err != nil {
		return nil, err
	}

	agentClient := agentv1.NewAgentServiceClient(conn)

	// Registration
	registrationClient := identityinfra.NewRegistrationClient(agentClient)
	registerClusterUC := identityuc.NewRegisterClusterUC(registrationClient)

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
	slog.Info("registration completed successfully!")

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

func StartTunnelLoop(client agentv1.AgentServiceClient, clusterID string) {
	backoff := time.Second * 1
	maxBackoff := time.Minute * 1

	for {
		log.Printf("Cluster %s: Attempting to connect to Control Plane...", clusterID)

		// 1. Setup context và metadata
		md := metadata.Pairs("x-cluster-id", clusterID)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		// 2. Open stream
		stream, err := client.ConnectTunnel(ctx)
		if err != nil {
			log.Printf("Connection failed: %v. Retrying in %v", err, backoff)
			time.Sleep(backoff)

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		backoff = time.Second * 1
		log.Println("Tunnel established successfully.")

		handleStream(stream, clusterID)

		log.Println("Stream connection lost. Reconnecting...")
	}
}

func handleStream(stream agentv1.AgentService_ConnectTunnelClient, clusterID string) {
	done := make(chan struct{})

	// Goroutine A: Listen cmd from server
	go func() {
		for {
			cmd, err := stream.Recv()
			if err != nil {
				log.Printf("Recv error: %v", err)
				close(done)
				return
			}
			log.Printf("Executing command: %s", cmd.Action)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Send heartbeat
			err := stream.Send(&agentv1.ConnectTunnelRequest{
				ClusterId: clusterID,
				Status:    "ONLINE",
			})
			if err != nil {
				log.Printf("Send error: %v", err)
				return
			}
		}
	}
}
