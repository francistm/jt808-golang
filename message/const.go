package message

import "errors"

const (
	AckType_0001_Succeed byte = iota
	AckType_0001_Failed
	AckType_0001_Incorrect
	AckType_0001_Unsupport
)

const (
	AckType_8001_Succeed byte = iota
	AckType_8001_Failed
	AckType_8001_Incorrect
	AckType_8001_Unsupport
	AckType_8001_AlarmConfirm
)

const (
	AckType_8100_Succeed byte = iota
	AckType_8100_Vehicle_Registered
	AckType_8100_Vehicle_NotFound
	AckType_8100_Terminal_Registered
	AckType_8100_Terminal_NotFound
)

const (
	Version2013 byte = iota
	Version2019
)

var ErrMesgNotSupport = errors.New("unsupported message")
