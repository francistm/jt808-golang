package message

// 终端注册
type Body0100 struct {
	Province          uint16
	City              uint16
	Manufacturer      []byte `jt808:"11,raw"`
	DeviceModel       []byte `jt808:"30,raw"`
	DeviceId          []byte `jt808:"30,raw"`
	LicencePlateColor uint8
	LicencePlate      string `jt808:"-1,gbk"`
}
