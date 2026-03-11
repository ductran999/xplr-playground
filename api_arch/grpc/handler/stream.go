package handler

import (
	"play-ground/api_arch/grpc/gen/proto/streampb"
	"time"
)

type WatchProvision struct {
	streampb.UnimplementedProvisionServiceServer
}

func (s *WatchProvision) WatchProvision(
	req *streampb.ProvisionRequest,
	stream streampb.ProvisionService_WatchProvisionServer,
) error {

	steps := []string{
		"Creating PVC",
		"Starting CNPG",
		"Deploying App",
		"Creating Service",
	}

	for _, step := range steps {
		err := stream.Send(&streampb.ProvisionEvent{
			Step:   step,
			Status: "IN_PROGRESS",
		})
		if err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}
