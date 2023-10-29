package message

// 平台通用应答
type Body8001 struct {
	AckMesgId   uint16
	AckSerialId uint16
	AckType     uint8
}
