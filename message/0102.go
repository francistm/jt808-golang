package message

// 终端鉴权
type Body0102 struct {
	AuthCodeSize  uint8
	AuthCode      string
	DeviceIMEI    []byte `jt808:"15,raw"`
	DeviceVersion []byte `jt808:"20,raw"`
}
