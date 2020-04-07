package jt808

import (
	"bytes"
)

type MessagePack struct {
	PackBody      PackBody
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool

	bodyBuf       []byte
	unmarshalFunc func([]byte) (PackBody, error)
}

func (ptr *MessagePack) ConcatAndUnmarshal(packs ...*MessagePack) error {
	buf := bytes.NewBuffer(ptr.bodyBuf)

	for _, pack := range packs {
		buf.Write(pack.bodyBuf)
	}

	if packBody, err := ptr.unmarshalFunc(buf.Bytes()); err != nil {
		return err
	} else {
		ptr.PackBody = packBody
	}

	// cleanup after unmarshal
	ptr.bodyBuf = nil
	ptr.unmarshalFunc = nil

	return nil
}

type PackBody interface {
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

	if body.AcknowledgeType, err = reader.ReadByte(); err != nil {
		return
	}

	return
}

type Body0200 struct {
}

func unmarshalBody0200(buf []byte) (body Body0200, err error) {
	reader := bytes.NewReader(buf)
	b := make([]byte, 28)

	// TODO 实现0200消息解析
	_, err = reader.Read(b)

	return
}

type Body0801 struct {
	MediaId          uint32
	MediaType        uint8
	MediaContentType uint8
	EventId          uint8
	ChannelId        uint8
	PackBody0200     Body0200
	MediaContent     []byte
}

func unmarshalBody0801(buf []byte) (body Body0801, err error) {
	reader := bytes.NewReader(buf)

	var mediaContentBuf bytes.Buffer

	if body.MediaId, err = readUint32(reader); err != nil {
		return
	}

	if body.MediaType, err = reader.ReadByte(); err != nil {
		return
	}

	if body.MediaContentType, err = reader.ReadByte(); err != nil {
		return
	}

	if body.EventId, err = reader.ReadByte(); err != nil {
		return
	}

	if body.ChannelId, err = reader.ReadByte(); err != nil {
		return
	}

	packBody0200Buf := make([]byte, 28)

	if _, err = reader.Read(packBody0200Buf); err != nil {
		return
	} else if body.PackBody0200, err = unmarshalBody0200(packBody0200Buf); err != nil {
		return
	}

	if _, err = reader.WriteTo(&mediaContentBuf); err != nil {
		return
	} else {
		body.MediaContent = mediaContentBuf.Bytes()
	}

	return
}
