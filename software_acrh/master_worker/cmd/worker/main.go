package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"
	"play-ground/software_acrh/master_worker/internal/worker/app"
	"play-ground/software_acrh/master_worker/internal/worker/config"

	"google.golang.org/grpc/metadata"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	w, err := app.Initialize(cfg)
	if err != nil {
		log.Fatalf("init worker failed: %v", err)
	}
	defer w.Close()

	if err := w.Run(appCtx); err != nil {
		log.Fatalf("agent start error: %v", err)
	}

	// StartTunnelLoop(client, "cluster-uuid-123")
}

func StartTunnelLoop(client pb.AgentServiceClient, clusterID string) {
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

func handleStream(stream pb.AgentService_ConnectTunnelClient, clusterID string) {
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
			err := stream.Send(&pb.ConnectTunnelRequest{
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
