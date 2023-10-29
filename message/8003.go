package message

// 终端注销
type Body8003 struct {
	SerialId   uint16
	TotalCount uint16
	RawMesgIds []byte `jt808:",raw"`
}
