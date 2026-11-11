package engine

import (
	"context"
	"testing"

	"github.com/databasifylabs/engine-sdk/contract"
	"github.com/databasifylabs/engine-sdk/enums"
)

func TestTheTemplateImplementsEveryGroup(t *testing.T) {
	var _ contract.Engine = New()
	if _, partial := New().Nodes().(contract.Partial); !partial {
		t.Fatal("the nodes stub must list what it implements")
	}
	if _, partial := New().Cluster().(contract.Partial); !partial {
		t.Fatal("the cluster stub must list what it implements")
	}
}

func TestStatusReportsTheRelease(t *testing.T) {
	t.Setenv("DATABASIFY_RELEASE", "9.1.3")
	status, err := New().Nodes().Status(context.Background())
	if err != nil || status.Role != enums.NodeRolePrimary || status.Version == nil || *status.Version != "9.1.3" {
		t.Fatalf("status = %+v, err = %v", status, err)
	}
}
