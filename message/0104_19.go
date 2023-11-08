package message

import "github.com/francistm/jt808-golang/internal/bytes"

type Body0104_19 struct {
	// 终端类型
	TerminalKind byte

	// 制造商ID
	// 5位
	ManufacturerId []byte

	// 终端型号
	// 30位
	TerminalModel []byte

	// 终端ID
	// 30位
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

func (b *Body0104_19) MarshalBinary() ([]byte, error) {
	buffer := bytes.NewBuffer()

	buffer.WriteByte(b.TerminalKind)
	buffer.Write(b.ManufacturerId)
	buffer.Write(b.TerminalModel)
	buffer.Write(b.TerminalId)
	buffer.WriteBCD(b.TerminalIccid, 10)
	buffer.WriteByte(uint8(len(b.TerminalHardwareVersion)))
	buffer.Write(b.TerminalHardwareVersion)
	buffer.WriteByte(uint8(len(b.TerminalFirmwareVersion)))
	buffer.Write(b.TerminalFirmwareVersion)
	buffer.WriteByte(b.GnssModuleProperty)
	buffer.WriteByte(b.CommunicationModuleProperty)

	return buffer.Bytes(), nil
}

func (b *Body0104_19) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)

	terminalKind, err := reader.ReadByte()

	if err != nil {
		return err
	}

	manufacturerId, err := reader.ReadFixedBytes(5)

	if err != nil {
		return err
	}

	terminalModel, err := reader.ReadFixedBytes(30)

	if err != nil {
		return err
	}

	terminalId, err := reader.ReadFixedBytes(30)

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

	b.TerminalKind = terminalKind
	b.ManufacturerId = manufacturerId
	b.TerminalModel = terminalModel
	b.TerminalId = terminalId
	b.TerminalIccid = terminalIccid
	b.TerminalHardwareVersionLength = terminalHardwareVersionLength
	b.TerminalHardwareVersion = terminalHardwareVersion
	b.TerminalFirmwareVersionLength = terminalFirmwareVersionLength
	b.TerminalFirmwareVersion = terminalFirmwareVersion
	b.GnssModuleProperty = gnssModuleProperty
	b.CommunicationModuleProperty = communicationModuleProperty

	return nil
}
