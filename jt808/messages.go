package jt808

import "bytes"

type MessagePack struct {
	PackBody   interface{}
	PackHeader PackHeader
	Checksum   uint8
}

type Body0001 struct {
	AcknowledgeSerialId  uint16
	AcknowledgeMessageId uint16
	AcknowledgeType      uint8
}

func unmarshalBody0001(buf []byte) (body Body0001, err error) {
	reader := bytes.NewReader(buf)

	if body.AcknowledgeMessageId, err = readUint16(reader); err != nil {
		return
	}

	if body.AcknowledgeSerialId, err = readUint16(reader); err != nil {
		return
	}

	if body.AcknowledgeType, err = readUint8(reader); err != nil {
		return
	}

	return
}
