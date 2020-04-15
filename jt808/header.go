package jt808

import (
	"bytes"
)

// 消息体头部
type PackHeader struct {
	// 消息 ID
	MessageId uint16

	// 消息体属性
	Property PackProperty

	// 终端手机号
	TerminalMobileNo string

	// 消息流水号
	SerialNo uint16

	// 消息包封装项
	Package *PackPackage
}

// 消息体属性
type PackProperty struct {
	// 消息体长度
	BodyByteLength uint16

	// 数据加密方式
	IsEncrypted bool

	// 是否为长消息进行分包
	IsMultiplePackage bool
}

// 消息包分装项
type PackPackage struct {
	// 消息总包数
	TotalCount uint16

	// 包序号，由 1 开始
	CurrentIndex uint16
}

// 由二进制解析一个消息包的头部
// 如果起始字节为 0x7e 则跳过该字节
func UnmarshalHeader(in []byte, v *PackHeader) error {
	reader := bytes.NewReader(in)

	if b, err := reader.ReadByte(); err != nil {
		return err
	} else if b != 0x7e {
		_ = reader.UnreadByte()
	}

	if i, err := readUint16(reader); err != nil {
		return err
	} else {
		v.MessageId = i
	}

	if i, err := readUint16(reader); err != nil {
		return err
	} else {
		v.Property.BodyByteLength = i & 0x03ff
		v.Property.IsEncrypted = ((i >> 10) & 0x01) == 0x01
		v.Property.IsMultiplePackage = ((i >> 13) & 0x01) == 0x01
	}

	if s, err := readBCD(reader, 6); err != nil {
		return err
	} else {
		v.TerminalMobileNo = s
	}

	if i, err := readUint16(reader); err != nil {
		return err
	} else {
		v.SerialNo = i
	}

	if v.Property.IsMultiplePackage {
		packPackagePtr := new(PackPackage)

		if i, err := readUint16(reader); err != nil {
			return err
		} else {
			packPackagePtr.TotalCount = i
		}

		if i, err := readUint16(reader); err != nil {
			return err
		} else {
			packPackagePtr.CurrentIndex = i
		}

		v.Package = packPackagePtr
	}

	return nil
}
