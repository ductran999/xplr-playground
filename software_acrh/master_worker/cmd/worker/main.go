package main

import (
	"context"
	"log"
	"time"

	pb "play-ground/software_acrh/master_worker/api/gen/pb/agent/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	StartTunnel(client, "cluster-uuid-123")
}

func StartTunnel(client pb.AgentServiceClient, clusterID string) {
	stream, err := client.ConnectTunnel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Goroutine 1: Receive command
	go func() {
		for {
			cmd, err := stream.Recv()
			if err != nil {
				log.Printf("Stream closed: %v", err)
				return
			}
			log.Printf("Received command: %s", cmd.Action)
		}
	}()

	// Goroutine 2: send Heartbeat
	for {
		log.Println("Agent ping server")
		stream.Send(&pb.ConnectTunnelRequest{
			ClusterId: clusterID,
			Status:    "HEALTHY",
		})
		time.Sleep(10 * time.Second)
	}
}
