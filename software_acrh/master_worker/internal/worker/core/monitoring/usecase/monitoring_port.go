package monitoringuc

import "context"

type StatusReporter interface {
	SendHeartbeat(ctx context.Context, agentID string, status string) error
}
