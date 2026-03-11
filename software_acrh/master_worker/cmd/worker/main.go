package main

import (
	"context"
	"log"
	"time"

	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAgentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &pb.RegisterRequest{
		RegistrationToken: "OK",
		Hostname:          "my-k8s-node-01",
		AgentVersion:      "v1.0.0",
	})

	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Registered successfully! ClusterID: %s", resp.ClusterId)

	StartTunnelLoop(client, "cluster-uuid-123")
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
