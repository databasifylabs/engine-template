package engine

import (
	"context"
	"os"

	"github.com/databasifylabs/engine-sdk/contract"
	"github.com/databasifylabs/engine-sdk/enums"
	"github.com/databasifylabs/engine-sdk/requests"
	"github.com/databasifylabs/engine-sdk/responses"
)

type Nodes struct{}

func (Nodes) Implements() []string {
	return []string{"node status"}
}

func (Nodes) Status(context.Context) (responses.NodeStatus, error) {
	release := os.Getenv("DATABASIFY_RELEASE")

	return responses.NodeStatus{Role: enums.NodeRolePrimary, Version: &release}, nil
}

func (Nodes) Configure(context.Context, requests.NodeConfigure) (responses.NodeConfigure, error) {
	return responses.NodeConfigure{}, contract.ErrUnsupported
}

func (Nodes) Start(context.Context) error { return contract.ErrUnsupported }

func (Nodes) Restart(context.Context) error { return contract.ErrUnsupported }

func (Nodes) Wake(context.Context, requests.NodeWake) error { return contract.ErrUnsupported }

func (Nodes) Ready(context.Context, requests.NodeReady, contract.Report) error {
	return contract.ErrUnsupported
}

func (Nodes) Bootstrap(context.Context, requests.NodeBootstrap) error {
	return contract.ErrUnsupported
}

func (Nodes) Metrics(context.Context) (responses.NodeMetrics, error) {
	return responses.NodeMetrics{}, contract.ErrUnsupported
}

func (Nodes) ApplyConfig(context.Context) error { return contract.ErrUnsupported }

func (Nodes) SyncConfig(context.Context) (responses.ConfigSync, error) {
	return responses.ConfigSync{}, contract.ErrUnsupported
}

func (Nodes) ReloadTLS(context.Context) error { return contract.ErrUnsupported }

func (Nodes) InstallTLS(context.Context, requests.TlsInstall) error {
	return contract.ErrUnsupported
}

func (Nodes) Reset(context.Context, requests.NodeReset) error { return contract.ErrUnsupported }

func (Nodes) Insights(context.Context) (responses.QueryInsights, error) {
	return responses.QueryInsights{}, contract.ErrUnsupported
}
