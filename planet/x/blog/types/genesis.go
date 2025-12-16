package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		PortId: PortID, PostList: []Post{}, SentPostList: []SentPost{}, TimeoutPostList: []TimeoutPost{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	postIdMap := make(map[uint64]bool)
	postCount := gs.GetPostCount()
	for _, elem := range gs.PostList {
		if _, ok := postIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for post")
		}
		if elem.Id >= postCount {
			return fmt.Errorf("post id should be lower or equal than the last id")
		}
		postIdMap[elem.Id] = true
	}
	sentPostIdMap := make(map[uint64]bool)
	sentPostCount := gs.GetSentPostCount()
	for _, elem := range gs.SentPostList {
		if _, ok := sentPostIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for sentPost")
		}
		if elem.Id >= sentPostCount {
			return fmt.Errorf("sentPost id should be lower or equal than the last id")
		}
		sentPostIdMap[elem.Id] = true
	}
	timeoutPostIdMap := make(map[uint64]bool)
	timeoutPostCount := gs.GetTimeoutPostCount()
	for _, elem := range gs.TimeoutPostList {
		if _, ok := timeoutPostIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for timeoutPost")
		}
		if elem.Id >= timeoutPostCount {
			return fmt.Errorf("timeoutPost id should be lower or equal than the last id")
		}
		timeoutPostIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
