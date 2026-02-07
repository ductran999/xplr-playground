package main

import (
	"context"
	"log"
	"time"

	"play-ground/grpc/proto/example.com/userpb"

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

	client := userpb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &userpb.GetUserRequest{
		Id: 2,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Id, resp.Name)
}
