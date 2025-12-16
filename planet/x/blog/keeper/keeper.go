package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"planet/x/blog/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	Port collections.Item[string]

	ibcKeeperFn    func() *ibckeeper.Keeper
	PostSeq        collections.Sequence
	Post           collections.Map[uint64, types.Post]
	SentPostSeq    collections.Sequence
	SentPost       collections.Map[uint64, types.SentPost]
	TimeoutPostSeq collections.Sequence
	TimeoutPost    collections.Map[uint64, types.TimeoutPost]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	ibcKeeperFn func() *ibckeeper.Keeper,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		ibcKeeperFn:    ibcKeeperFn,
		Port:           collections.NewItem(sb, types.PortKey, "port", collections.StringValue),
		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Post:           collections.NewMap(sb, types.PostKey, "post", collections.Uint64Key, codec.CollValue[types.Post](cdc)),
		PostSeq:        collections.NewSequence(sb, types.PostCountKey, "postSequence"),
		SentPost:       collections.NewMap(sb, types.SentPostKey, "sentPost", collections.Uint64Key, codec.CollValue[types.SentPost](cdc)),
		SentPostSeq:    collections.NewSequence(sb, types.SentPostCountKey, "sentPostSequence"),
		TimeoutPost:    collections.NewMap(sb, types.TimeoutPostKey, "timeoutPost", collections.Uint64Key, codec.CollValue[types.TimeoutPost](cdc)),
		TimeoutPostSeq: collections.NewSequence(sb, types.TimeoutPostCountKey, "timeoutPostSequence"),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
