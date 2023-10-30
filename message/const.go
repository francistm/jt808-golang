package message

const (
	AckTypeSucceed byte = iota
	AckTypeFailed
	AckTypeIncorrect
	AckTypeUnsupport
	AckTypeAlarmConfirm
)

const (
	Version2013 byte = iota
	Version2019
)
