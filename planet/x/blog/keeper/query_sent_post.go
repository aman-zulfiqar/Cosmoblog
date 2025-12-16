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

func (q queryServer) ListSentPost(ctx context.Context, req *types.QueryAllSentPostRequest) (*types.QueryAllSentPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sentPosts, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.SentPost,
		req.Pagination,
		func(_ uint64, value types.SentPost) (types.SentPost, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSentPostResponse{SentPost: sentPosts, Pagination: pageRes}, nil
}

func (q queryServer) GetSentPost(ctx context.Context, req *types.QueryGetSentPostRequest) (*types.QueryGetSentPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sentPost, err := q.k.SentPost.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetSentPostResponse{SentPost: sentPost}, nil
}
