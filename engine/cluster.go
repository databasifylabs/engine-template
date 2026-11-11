package engine

import (
	"context"

	"github.com/databasifylabs/engine-sdk/contract"
	"github.com/databasifylabs/engine-sdk/requests"
	"github.com/databasifylabs/engine-sdk/responses"
)

type Cluster struct{}

func (Cluster) Implements() []string {
	return []string{}
}

func (Cluster) Enable(context.Context, requests.ClusterEnable, contract.Report) error {
	return contract.ErrUnsupported
}

func (Cluster) Disable(context.Context) (responses.ClusterDisable, error) {
	return responses.ClusterDisable{}, contract.ErrUnsupported
}

func (Cluster) Failover(context.Context, requests.ClusterFailover) error {
	return contract.ErrUnsupported
}

func (Cluster) ForceQuorum(context.Context, requests.ClusterForceQuorum) (responses.ClusterForceQuorum, error) {
	return responses.ClusterForceQuorum{}, contract.ErrUnsupported
}

func (Cluster) Quorum(context.Context) (responses.ClusterQuorum, error) {
	return responses.ClusterQuorum{}, contract.ErrUnsupported
}

func (Cluster) Position(context.Context) (responses.ClusterPosition, error) {
	return responses.ClusterPosition{}, contract.ErrUnsupported
}

func (Cluster) Admit(context.Context, requests.MemberAdmit) error { return contract.ErrUnsupported }

func (Cluster) Promote(context.Context, requests.MemberPromote) error {
	return contract.ErrUnsupported
}

func (Cluster) Forget(context.Context, requests.MemberForget) error {
	return contract.ErrUnsupported
}

func (Cluster) Rejoin(context.Context, requests.ClusterRejoin, contract.Report) error {
	return contract.ErrUnsupported
}

func (Cluster) UpdateMembership(context.Context, requests.ClusterMembership) error {
	return contract.ErrUnsupported
}

func (Cluster) Watch(context.Context) error { return contract.ErrUnsupported }
