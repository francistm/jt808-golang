package message

const (
	AckTypeSucceed byte = iota
	AckTypeFailed
	AckTypeIncorrect
	AckTypeUnsupport
	AckTypeAlarmConfirm
)
