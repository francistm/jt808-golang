package message

import (
	"encoding/binary"

	"github.com/francistm/jt808-golang/internal/bytes"
)

// PackHeader 消息体头部
type PackHeader struct {
	// 消息 ID
	MessageID uint16

	// 消息体属性
	Property PackHeaderProperty

	// 协议版本号 (2019)
	VersionStep uint8

	// 终端手机号
	TerminalMobileNum string

	// 消息流水号
	SerialNum uint16

	// 消息包封装项
	Package *PackHeaderPackage
}

func (h *PackHeader) MarshalBinary() ([]byte, error) {
	var terminalMobileNumSize int

	if h.Property.Version == Version2013 {
		terminalMobileNumSize = 6
	} else if h.Property.Version == Version2019 {
		terminalMobileNumSize = 10
	}

	buf := bytes.NewBuffer()

	if err := buf.WriteUint16(h.MessageID); err != nil {
		return nil, err
	}

	b, err := h.Property.MarshalBinary()

	if err != nil {
		return nil, err
	}

	if _, err = buf.Write(b); err != nil {
		return nil, err
	}

	if err := buf.WriteBCD(h.TerminalMobileNum, terminalMobileNumSize); err != nil {
		return nil, err
	}

	if err := buf.WriteUint16(h.SerialNum); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *PackHeader) UnmarshalBinary(in []byte) error {
	var (
		err    error
		reader = bytes.NewReader(in)

		mesgId            uint16
		versionStep       uint8
		terminalMobileNum string
		serialNum         uint16
		packPackage       *PackHeaderPackage
	)

	mesgId, err = reader.ReadUint16()

	if err != nil {
		return err
	}

	propertyData, err := reader.ReadFixedBytes(2)

	if err != nil {
		return err
	}

	if err := h.Property.UnmarshalBinary(propertyData); err != nil {
		return err
	}

	if h.Property.Version == Version2013 {
		terminalMobileNum, err = reader.ReadBCD(6)

		if err != nil {
			return err
		}
	} else if h.Property.Version == Version2019 {
		versionStep, err = reader.ReadByte()

		if err != nil {
			return err
		}

		terminalMobileNum, err = reader.ReadBCD(10)

		if err != nil {
			return err
		}
	}

	serialNum, err = reader.ReadUint16()

	if err != nil {
		return err
	}

	if h.Property.IsMultiplePackage {
		pkgData, err := reader.ReadFixedBytes(4)

		if err != nil {
			return err
		}

		packPackage = new(PackHeaderPackage)

		if err := packPackage.UnmarshalBinary(pkgData); err != nil {
			return err
		}
	}

	h.MessageID = mesgId
	h.VersionStep = versionStep
	h.TerminalMobileNum = terminalMobileNum
	h.SerialNum = serialNum
	h.Package = packPackage

	return nil
}

type PackHeaderProperty struct {
	// 消息体长度
	BodyByteLength uint16

	// 数据加密方式
	IsEncrypted bool

	// 是否为长消息进行分包
	IsMultiplePackage bool

	// 版本表示
	Version uint8
}

func (p *PackHeaderProperty) MarshalBinary() ([]byte, error) {
	var (
		i   uint16
		out = make([]byte, 2)
	)

	i |= p.BodyByteLength

	if p.IsEncrypted {
		i |= 0x01 << 10
	}

	if p.IsMultiplePackage {
		i |= 0x01 << 13
	}

	if p.Version == Version2019 {
		i |= 0x01 << 14
	}

	binary.BigEndian.PutUint16(out, i)

	return out, nil
}

func (p *PackHeaderProperty) UnmarshalBinary(data []byte) error {
	i := binary.BigEndian.Uint16(data)

	p.BodyByteLength = i & 0x03ff
	p.IsEncrypted = ((i >> 10) & 0x01) == 0x01
	p.IsMultiplePackage = ((i >> 13) & 0x01) == 0x01
	p.Version = uint8(i >> 14 & 0x01)

	return nil
}

type PackHeaderPackage struct {
	// 消息总包数
	Total uint16

	// 包序号，由 1 开始
	Index uint16
}

func (p *PackHeaderPackage) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer()

	if err := b.WriteUint16(p.Total); err != nil {
		return nil, err
	}

	if err := b.WriteUint16(p.Index); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (p *PackHeaderPackage) UnmarshalBinary(b []byte) error {
	r := bytes.NewReader(b)

	total, err := r.ReadUint16()

	if err != nil {
		return err
	}

	index, err := r.ReadUint16()

	if err != nil {
		return err
	}

	p.Index = index
	p.Total = total

	return nil
}
