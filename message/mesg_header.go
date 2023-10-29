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
	Property PackProperty

	// 终端手机号
	TerminalMobileNo string

	// 消息流水号
	SerialNo uint16

	// 消息包封装项
	Package *PackPackage
}

func (h *PackHeader) MarshalBinary() ([]byte, error) {
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

	if err := buf.WriteBCD(h.TerminalMobileNo); err != nil {
		return nil, err
	}

	if err := buf.WriteUint16(h.SerialNo); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *PackHeader) UnmarshalBinary(in []byte) error {
	reader := bytes.NewReader(in)

	b, err := reader.ReadByte()

	if err != nil {
		return err
	}

	if b != 0x7e {
		_ = reader.UnreadByte()
	}

	i, err := reader.ReadUint16()

	if err != nil {
		return err
	}

	h.MessageID = i

	propertyData, err := reader.ReadBytes(2)

	if err != nil {
		return err
	}

	if err := h.Property.UnmarshalBinary(propertyData); err != nil {
		return err
	}

	s, err := reader.ReadBCD(6)

	if err != nil {
		return err
	}

	h.TerminalMobileNo = s

	i, err = reader.ReadUint16()

	if err != nil {
		return err
	}

	h.SerialNo = i

	if h.Property.IsMultiplePackage {
		pkgData := make([]byte, 4)

		if _, err := reader.Read(pkgData); err != nil {
			return err
		}

		h.Package = new(PackPackage)

		if err := h.Package.UnmarshalBinary(pkgData); err != nil {
			return err
		}
	}

	return nil
}

// PackProperty 消息体属性
type PackProperty struct {
	// 消息体长度
	BodyByteLength uint16

	// 数据加密方式
	IsEncrypted bool

	// 是否为长消息进行分包
	IsMultiplePackage bool
}

func (p *PackProperty) MarshalBinary() ([]byte, error) {
	var i uint16
	out := make([]byte, 2)

	i |= p.BodyByteLength

	if p.IsEncrypted {
		i |= 0x01 << 10
	}

	if p.IsMultiplePackage {
		i |= 0x01 << 13
	}

	binary.BigEndian.PutUint16(out, i)

	return out, nil
}

func (p *PackProperty) UnmarshalBinary(data []byte) error {
	i := binary.BigEndian.Uint16(data)

	p.BodyByteLength = i & 0x03ff
	p.IsEncrypted = ((i >> 10) & 0x01) == 0x01
	p.IsMultiplePackage = ((i >> 13) & 0x01) == 0x01

	return nil
}

// PackPackage 消息包分装项
type PackPackage struct {
	// 消息总包数
	TotalCount uint16

	// 包序号，由 1 开始
	Index uint16
}

func (p *PackPackage) MarshalBinary() ([]byte, error) {
	panic("not implemeneted")
}

func (p *PackPackage) UnmarshalBinary(b []byte) error {
	r := bytes.NewReader(b)

	totalCount, err := r.ReadUint16()

	if err != nil {
		return err
	}

	p.TotalCount = totalCount

	currentIndex, err := r.ReadUint16()

	if err != nil {
		return err
	}

	p.Index = currentIndex

	return nil
}
