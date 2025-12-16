package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"planet/x/blog/keeper"
	"planet/x/blog/types"
)

func createNTimeoutPost(keeper keeper.Keeper, ctx context.Context, n int) []types.TimeoutPost {
	items := make([]types.TimeoutPost, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		items[i].Title = strconv.Itoa(i)
		items[i].Chain = strconv.Itoa(i)
		items[i].Creator = strconv.Itoa(i)
		_ = keeper.TimeoutPost.Set(ctx, iu, items[i])
		_ = keeper.TimeoutPostSeq.Set(ctx, iu)
	}
	return items
}

func TestTimeoutPostQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNTimeoutPost(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetTimeoutPostRequest
		response *types.QueryGetTimeoutPostResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTimeoutPostRequest{Id: msgs[0].Id},
			response: &types.QueryGetTimeoutPostResponse{TimeoutPost: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetTimeoutPostRequest{Id: msgs[1].Id},
			response: &types.QueryGetTimeoutPostResponse{TimeoutPost: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetTimeoutPostRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetTimeoutPost(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestTimeoutPostQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNTimeoutPost(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTimeoutPostRequest {
		return &types.QueryAllTimeoutPostRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListTimeoutPost(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.TimeoutPost), step)
			require.Subset(t, msgs, resp.TimeoutPost)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListTimeoutPost(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.TimeoutPost), step)
			require.Subset(t, msgs, resp.TimeoutPost)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListTimeoutPost(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.TimeoutPost)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListTimeoutPost(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
