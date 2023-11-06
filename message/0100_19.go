package message

// 终端注册
type Body0100_19 struct {
	Province          uint16
	City              uint16
	Manufacturer      string `jt808:"11,raw"`
	DeviceModel       string `jt808:"30,raw"`
	DeviceId          string `jt808:"30,raw"`
	LicencePlateColor uint8
	LicencePlate      string `jt808:"-1,gbk"`
}
