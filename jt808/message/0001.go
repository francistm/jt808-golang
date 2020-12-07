package message

// Body0001 0x0001 消息体正文结构体
type Body0001 struct {
	AcknowledgeMessageID uint16 `jt808:""`
	AcknowledgeSerialID  uint16 `jt808:""`
	AcknowledgeType      uint8  `jt808:""`
}
