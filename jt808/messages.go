package jt808

import (
	"bytes"
	"errors"
	"fmt"
)

// MessagePack 通用的消息体结构
type MessagePack struct {
	PackBody      interface{}
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool

	bodyBuf []byte
}

// ConcatAndUnmarshal 拼接多个分段消息并解析
func (messagePack *MessagePack) ConcatAndUnmarshal(packs ...*MessagePack) error {
	buf := bytes.NewBuffer(messagePack.bodyBuf)

	if messagePack.PackHeader.Package == nil {
		return errors.New("cannot concat packages without package property header")
	}

	concatPackageIndexList := make(map[uint16]bool)
	concatPackageIndexList[messagePack.PackHeader.Package.CurrentIndex] = true

	for _, pack := range packs {
		if pack.PackHeader.Package == nil {
			return errors.New("package to be concat doesn't have package property header")
		}

		if !concatPackageIndexList[pack.PackHeader.Package.CurrentIndex] {
			buf.Write(pack.bodyBuf)
			concatPackageIndexList[pack.PackHeader.Package.CurrentIndex] = true
		}
	}

	for i := uint16(0); i < messagePack.PackHeader.Package.TotalCount; i++ {
		if concatPackageIndexList[i+1] == false {
			return fmt.Errorf("missing package with index %d to concat and unmarshal message", i+1)
		}
	}

	return nil
}
