package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"play-ground/api_arch/grpc/gen/proto/streampb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := streampb.NewProvisionServiceClient(conn)

	stream, _ := client.WatchProvision(
		context.Background(),
		&streampb.ProvisionRequest{
			WorkspaceId: "ws-123",
		},
	)

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(msg.Step, msg.Status)
	}

}
