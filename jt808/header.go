package jt808

import (
	"bytes"
	"encoding/binary"
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

// PackProperty 消息体属性
type PackProperty struct {
	// 消息体长度
	BodyByteLength uint16

	// 数据加密方式
	IsEncrypted bool

	// 是否为长消息进行分包
	IsMultiplePackage bool
}

func (ptr *PackProperty) marshal() ([]byte, error) {
	var i uint16
	out := make([]byte, 2)

	i |= ptr.BodyByteLength

	if ptr.IsEncrypted {
		i |= 0x01 << 10
	}

	if ptr.IsMultiplePackage {
		i |= 0x01 << 13
	}

	binary.BigEndian.PutUint16(out, i)

	return out, nil
}

// PackPackage 消息包分装项
type PackPackage struct {
	// 消息总包数
	TotalCount uint16

	// 包序号，由 1 开始
	CurrentIndex uint16
}

func marshalHeader(v *PackHeader) ([]byte, error) {
	var buf bytes.Buffer

	if err := writeUint16(v.MessageID, &buf); err != nil {
		return nil, err
	}

	if b, err := v.Property.marshal(); err != nil {
		return nil, err
	} else if _, err = buf.Write(b); err != nil {
		return nil, err
	}

	if err := writeBCD(v.TerminalMobileNo, &buf); err != nil {
		return nil, err
	}

	if err := writeUint16(v.SerialNo, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalHeader 由二进制解析一个消息包的头部
// 如果起始字节为 0x7e 则跳过该字节
func UnmarshalHeader(in []byte, v *PackHeader) error {
	reader := bytes.NewReader(in)

	if b, err := reader.ReadByte(); err != nil {
		return err
	} else if b != 0x7e {
		_ = reader.UnreadByte()
	}

	i, err := readUint16(reader)

	if err != nil {
		return err
	}

	v.MessageID = i

	i, err = readUint16(reader)

	if err != nil {
		return err
	}

	v.Property.BodyByteLength = i & 0x03ff
	v.Property.IsEncrypted = ((i >> 10) & 0x01) == 0x01
	v.Property.IsMultiplePackage = ((i >> 13) & 0x01) == 0x01

	s, err := readBCD(reader, 6)

	if err != nil {
		return err
	}

	v.TerminalMobileNo = s

	i, err = readUint16(reader)

	if err != nil {
		return err
	}

	v.SerialNo = i

	if v.Property.IsMultiplePackage {
		packPackagePtr := new(PackPackage)

		i, err := readUint16(reader)

		if err != nil {
			return err
		}
		packPackagePtr.TotalCount = i

		i, err = readUint16(reader)

		if err != nil {
			return err
		}

		packPackagePtr.CurrentIndex = i

		v.Package = packPackagePtr
	}

	return nil
}
