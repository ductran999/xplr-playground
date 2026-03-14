package gclient

import (
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTimeout = 10 * time.Second
)

var (
	ErrMissingAddress = errors.New("gRPC address is required")
)

type Config struct {
	Address string
	Timeout time.Duration
}

func (cfg *Config) validateAndSetDefaults() error {
	if cfg.Address == "" {
		return ErrMissingAddress
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}

	return nil
}

func Connect(cfg Config) (*grpc.ClientConn, error) {
	if err := cfg.validateAndSetDefaults(); err != nil {
		return nil, err
	}

	return grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
