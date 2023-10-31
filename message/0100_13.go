package message

// 终端注册
type Body0100_13 struct {
	Province          uint16
	City              uint16
	Manufacturer      []byte `jt808:"5,raw"`
	DeviceModel       []byte `jt808:"20,raw"`
	DeviceId          []byte `jt808:"7,raw"`
	LicencePlateColor uint8
	LicencePlate      string `jt808:"-1,gbk"`
}
