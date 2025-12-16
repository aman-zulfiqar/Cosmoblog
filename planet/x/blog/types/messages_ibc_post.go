package types

func NewMsgSendIbcPost(
	creator string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	title string,
	content string,
) *MsgSendIbcPost {
	return &MsgSendIbcPost{
		Creator:          creator,
		Port:             port,
		ChannelID:        channelID,
		TimeoutTimestamp: timeoutTimestamp,
		Title:            title,
		Content:          content,
	}
}
