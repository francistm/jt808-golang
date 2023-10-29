package message

// 终端注册应答
type Body8100 struct {
	SerialId uint16
	Result   uint8
	AuthCode string
}
