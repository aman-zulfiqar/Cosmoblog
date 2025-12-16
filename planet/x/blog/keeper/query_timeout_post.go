package keeper

import (
	"context"
	"errors"

	"planet/x/blog/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListTimeoutPost(ctx context.Context, req *types.QueryAllTimeoutPostRequest) (*types.QueryAllTimeoutPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	timeoutPosts, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.TimeoutPost,
		req.Pagination,
		func(_ uint64, value types.TimeoutPost) (types.TimeoutPost, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTimeoutPostResponse{TimeoutPost: timeoutPosts, Pagination: pageRes}, nil
}

func (q queryServer) GetTimeoutPost(ctx context.Context, req *types.QueryGetTimeoutPostRequest) (*types.QueryGetTimeoutPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	timeoutPost, err := q.k.TimeoutPost.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetTimeoutPostResponse{TimeoutPost: timeoutPost}, nil
}
