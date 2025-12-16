package keeper_test

import (
	"testing"

	"planet/x/blog/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:           types.DefaultParams(),
		PortId:           types.PortID,
		PostList:         []types.Post{{Id: 0}, {Id: 1}},
		PostCount:        2,
		SentPostList:     []types.SentPost{{Id: 0}, {Id: 1}},
		SentPostCount:    2,
		TimeoutPostList:  []types.TimeoutPost{{Id: 0}, {Id: 1}},
		TimeoutPostCount: 2,
	}
	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, genesisState.PortId, got.PortId)
	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.PostList, got.PostList)
	require.Equal(t, genesisState.PostCount, got.PostCount)
	require.EqualExportedValues(t, genesisState.SentPostList, got.SentPostList)
	require.Equal(t, genesisState.SentPostCount, got.SentPostCount)
	require.EqualExportedValues(t, genesisState.TimeoutPostList, got.TimeoutPostList)
	require.Equal(t, genesisState.TimeoutPostCount, got.TimeoutPostCount)

}
