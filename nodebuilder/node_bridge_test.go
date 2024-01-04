package nodebuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sunrise-zone/sunrise-node/core"
	coremodule "github.com/sunrise-zone/sunrise-node/nodebuilder/core"
	"github.com/sunrise-zone/sunrise-node/nodebuilder/node"
	"github.com/sunrise-zone/sunrise-node/nodebuilder/p2p"
)

func TestBridge_WithMockedCoreClient(t *testing.T) {
	t.Skip("skipping") // consult https://github.com/celestiaorg/celestia-core/issues/667 for reasoning
	repo := MockStore(t, DefaultConfig(node.Bridge))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	client := core.StartTestNode(t).Client
	node, err := New(node.Bridge, p2p.Private, repo, coremodule.WithClient(client))
	require.NoError(t, err)
	require.NotNil(t, node)
	err = node.Start(ctx)
	require.NoError(t, err)

	err = node.Stop(ctx)
	require.NoError(t, err)
}
