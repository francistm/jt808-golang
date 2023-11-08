package message

import "github.com/francistm/jt808-golang/internal/bytes"

// 查询终端属性应答
type Body0107_13 struct {
	// 终端类型
	TerminalKind byte

	// 制造商ID
	// 5位
	ManufacturerId []byte

	// 终端型号
	// 20位
	TerminalModel []byte

	// 终端ID
	// 7位
	TerminalId []byte

	// 终端SIM卡ICCID
	// 10位 BCD
	TerminalIccid string

	// 终端硬件版本号长度
	TerminalHardwareVersionLength byte

	// 终端硬件版本号
	TerminalHardwareVersion []byte

	// 终端固件版本号长度
	TerminalFirmwareVersionLength byte

	// 终端固件版本号
	TerminalFirmwareVersion []byte

	// GNSS模块属性
	GnssModuleProperty byte

	// 通信模块属性
	CommunicationModuleProperty byte
}

func (v *Body0107_13) MarshalBinary() ([]byte, error) {
	buffer := bytes.NewBuffer()

	buffer.WriteByte(v.TerminalKind)
	buffer.Write(v.ManufacturerId)
	buffer.Write(v.TerminalModel)
	buffer.Write(v.TerminalId)
	buffer.WriteBCD(v.TerminalIccid, 10)
	buffer.WriteByte(uint8(len(v.TerminalHardwareVersion)))
	buffer.Write(v.TerminalHardwareVersion)
	buffer.WriteByte(uint8(len(v.TerminalFirmwareVersion)))
	buffer.Write(v.TerminalFirmwareVersion)
	buffer.WriteByte(v.GnssModuleProperty)
	buffer.WriteByte(v.CommunicationModuleProperty)

	return buffer.Bytes(), nil
}

func (v *Body0107_13) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)

	terminalKind, err := reader.ReadByte()

	if err != nil {
		return err
	}

	manufacturerId, err := reader.ReadFixedBytes(5)

	if err != nil {
		return err
	}

	terminalModel, err := reader.ReadFixedBytes(20)

	if err != nil {
		return err
	}

	terminalId, err := reader.ReadFixedBytes(7)

	if err != nil {
		return err
	}

	terminalIccid, err := reader.ReadBCD(10)

	if err != nil {
		return err
	}

	terminalHardwareVersionLength, err := reader.ReadByte()

	if err != nil {
		return err
	}

	terminalHardwareVersion, err := reader.ReadFixedBytes(int(terminalHardwareVersionLength))

	if err != nil {
		return err
	}

	terminalFirmwareVersionLength, err := reader.ReadByte()

	if err != nil {
		return err
	}

	terminalFirmwareVersion, err := reader.ReadFixedBytes(int(terminalFirmwareVersionLength))

	if err != nil {
		return err
	}

	gnssModuleProperty, err := reader.ReadByte()

	if err != nil {
		return err
	}

	communicationModuleProperty, err := reader.ReadByte()

	if err != nil {
		return err
	}

	v.TerminalKind = terminalKind
	v.ManufacturerId = manufacturerId
	v.TerminalModel = terminalModel
	v.TerminalId = terminalId
	v.TerminalIccid = terminalIccid
	v.TerminalHardwareVersionLength = terminalHardwareVersionLength
	v.TerminalHardwareVersion = terminalHardwareVersion
	v.TerminalFirmwareVersionLength = terminalFirmwareVersionLength
	v.TerminalFirmwareVersion = terminalFirmwareVersion
	v.GnssModuleProperty = gnssModuleProperty
	v.CommunicationModuleProperty = communicationModuleProperty

	return nil
}
