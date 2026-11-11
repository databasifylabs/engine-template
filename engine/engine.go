package engine

import "github.com/databasifylabs/engine-sdk/contract"

type Engine struct{}

func New() Engine {
	return Engine{}
}

func (Engine) Name() string { return "example" }

func (Engine) Capabilities() contract.Capabilities { return contract.Capabilities{} }

func (Engine) HostProfile() contract.HostProfile { return contract.HostProfile{} }

func (Engine) Nodes() contract.Nodes { return Nodes{} }

func (Engine) Cluster() contract.Cluster { return Cluster{} }

func (Engine) Accounts() contract.Accounts { return contract.UnsupportedAccounts{} }

func (Engine) Archive() contract.Archive { return contract.UnsupportedArchive{} }

func (Engine) Migration() contract.Migration { return contract.UnsupportedMigration{} }

func (Engine) Pools() contract.Pools { return contract.UnsupportedPools{} }

func (Engine) Config() contract.Config { return contract.UnsupportedConfig{} }
