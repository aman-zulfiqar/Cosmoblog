package blog

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"planet/x/blog/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListPost",
					Use:       "list-post",
					Short:     "List all post",
				},
				{
					RpcMethod:      "GetPost",
					Use:            "get-post [id]",
					Short:          "Gets a post by id",
					Alias:          []string{"show-post"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ListSentPost",
					Use:       "list-sent-post",
					Short:     "List all sentPost",
				},
				{
					RpcMethod:      "GetSentPost",
					Use:            "get-sent-post [id]",
					Short:          "Gets a sentPost by id",
					Alias:          []string{"show-sent-post"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "ListTimeoutPost",
					Use:       "list-timeout-post",
					Short:     "List all timeoutPost",
				},
				{
					RpcMethod:      "GetTimeoutPost",
					Use:            "get-timeout-post [id]",
					Short:          "Gets a timeoutPost by id",
					Alias:          []string{"show-timeout-post"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
