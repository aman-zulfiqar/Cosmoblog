package types

// GetBytes is a helper for serialising
func (p IbcPostPacketData) GetBytes() ([]byte, error) {
	var modulePacket BlogPacketData

	modulePacket.Packet = &BlogPacketData_IbcPostPacket{&p}

	return modulePacket.Marshal()
}
