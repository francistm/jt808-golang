package message

// 终端通用应答
type Body0001 struct {
	AckMesgId   uint16
	AckSerialId uint16
	AckType     uint8
}
