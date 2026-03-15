package monitoringuc

import (
	"context"
)

type MonitoringUseCase interface {
	SendHeartbeat(ctx context.Context, agentID string) error
	SendClusterStats(ctx context.Context) error
}

type monitoringUC struct {
	statusReporter StatusReporter
}

func NewMonitoringUseCase(statusReporter StatusReporter) MonitoringUseCase {
	return &monitoringUC{
		statusReporter: statusReporter,
	}
}

func (uc *monitoringUC) SendHeartbeat(ctx context.Context, agentID string) error {
	status := "HEALTHY"
	if err := uc.statusReporter.SendHeartbeat(ctx, agentID, status); err != nil {
		return err
	}

	return nil
}

func (uc *monitoringUC) SendClusterStats(ctx context.Context) error {
	return nil
}
