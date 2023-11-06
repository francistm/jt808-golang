package message

// 平台通用应答
type Body8001 struct {
	AckSerialId uint16
	AckMesgId   uint16
	AckType     uint8
}
