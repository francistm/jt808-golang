package jt808

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"time"
)

var locationCST = time.FixedZone("Asia/Shanghai", 8*3600)

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

	if ptr.PackHeader.Package == nil {
		return errors.New("cannot concat packages without package property header")
	}

	concatPackageIndexList := make(map[uint16]bool)
	concatPackageIndexList[ptr.PackHeader.Package.CurrentIndex] = true

	for _, pack := range packs {
		if pack.PackHeader.Package == nil {
			return errors.New("package to be concat doesn't have package property header")
		}

		if !concatPackageIndexList[pack.PackHeader.Package.CurrentIndex] {
			buf.Write(pack.bodyBuf)
			concatPackageIndexList[pack.PackHeader.Package.CurrentIndex] = true
		}
	}

	for i := uint16(0); i < ptr.PackHeader.Package.TotalCount; i++ {
		if concatPackageIndexList[i+1] == false {
			return errors.New(fmt.Sprintf("missing package with index %d to concat and unmarshal message", i+1))
		}
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
	WarnFlag     uint32
	StatusFlag   uint32
	Latitude     float64
	Longitude    float64
	Altitude     uint16
	Speed        float32
	Direction    uint16
	Time         time.Time
	ExtraMessage map[uint8][]byte
}

func unmarshalBody0200(buf []byte) (body Body0200, err error) {
	reader := bytes.NewReader(buf)

	if body.WarnFlag, err = readUint32(reader); err != nil {
		return
	}

	if body.StatusFlag, err = readUint32(reader); err != nil {
		return
	}

	if i, err := readUint32(reader); err != nil {
		return body, err
	} else {
		body.Latitude = float64(i) / math.Pow(10, 6)
	}

	if i, err := readUint32(reader); err != nil {
		return body, err
	} else {
		body.Longitude = float64(i) / math.Pow(10, 6)
	}

	if body.Altitude, err = readUint16(reader); err != nil {
		return
	}

	if i, err := readUint16(reader); err != nil {
		return body, err
	} else {
		body.Speed = float32(i) / 10
	}

	if body.Direction, err = readUint16(reader); err != nil {
		return
	}

	if s, err := readBCD(reader, 6); err != nil {
		return body, err
	} else {
		t, err := time.ParseInLocation("060102150405", s, locationCST)

		if err != nil {
			return body, err
		}

		body.Time = t
	}

	// read extra messages in 0200 (if have)
	for {
		extraDataId, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		if body.ExtraMessage == nil {
			body.ExtraMessage = make(map[uint8][]byte)
		}

		extraDataLength, err := reader.ReadByte()

		if err != nil {
			return body, err
		}

		body.ExtraMessage[extraDataId] = make([]byte, extraDataLength)

		if _, err := reader.Read(body.ExtraMessage[extraDataId]); err != nil {
			return body, err
		}
	}

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
