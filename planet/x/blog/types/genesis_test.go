package types_test

import (
	"testing"

	"planet/x/blog/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				PortId:   types.PortID,
				PostList: []types.Post{{Id: 0}, {Id: 1}}, PostCount: 2,
				SentPostList: []types.SentPost{{Id: 0}, {Id: 1}}, SentPostCount: 2, TimeoutPostList: []types.TimeoutPost{{Id: 0}, {Id: 1}}, TimeoutPostCount: 2}, valid: true,
		}, {
			desc: "duplicated post",
			genState: &types.GenesisState{
				PostList: []types.Post{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
				SentPostList: []types.SentPost{{Id: 0}, {Id: 1}}, SentPostCount: 2,
				TimeoutPostList: []types.TimeoutPost{{Id: 0}, {Id: 1}}, TimeoutPostCount: 2}, valid: false,
		}, {
			desc: "invalid post count",
			genState: &types.GenesisState{
				PostList: []types.Post{
					{
						Id: 1,
					},
				},
				PostCount:    0,
				SentPostList: []types.SentPost{{Id: 0}, {Id: 1}}, SentPostCount: 2,
				TimeoutPostList: []types.TimeoutPost{{Id: 0}, {Id: 1}}, TimeoutPostCount: 2}, valid: false,
		}, {
			desc: "duplicated sentPost",
			genState: &types.GenesisState{
				SentPostList: []types.SentPost{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
				TimeoutPostList: []types.TimeoutPost{{Id: 0}, {Id: 1}}, TimeoutPostCount: 2,
			}, valid: false,
		}, {
			desc: "invalid sentPost count",
			genState: &types.GenesisState{
				SentPostList: []types.SentPost{
					{
						Id: 1,
					},
				},
				SentPostCount:   0,
				TimeoutPostList: []types.TimeoutPost{{Id: 0}, {Id: 1}}, TimeoutPostCount: 2,
			}, valid: false,
		}, {
			desc: "duplicated timeoutPost",
			genState: &types.GenesisState{
				TimeoutPostList: []types.TimeoutPost{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		}, {
			desc: "invalid timeoutPost count",
			genState: &types.GenesisState{
				TimeoutPostList: []types.TimeoutPost{
					{
						Id: 1,
					},
				},
				TimeoutPostCount: 0,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
