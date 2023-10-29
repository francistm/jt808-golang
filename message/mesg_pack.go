package message

import (
	"fmt"

	"github.com/francistm/jt808-golang/internal"
	"github.com/francistm/jt808-golang/internal/bytes"
	"github.com/francistm/jt808-golang/internal/decode"
	"github.com/francistm/jt808-golang/internal/encode"
)

// MessagePack 通用的消息体结构
type MessagePack[T any] struct {
	PackBody      T
	PackHeader    PackHeader
	Checksum      uint8
	ChecksumValid bool
}

func (p *MessagePack[T]) MarshalBinary() ([]byte, error) {
	var (
		buf      = bytes.NewBuffer()
		bodyBuf  = bytes.NewBuffer()
		finalBuf = bytes.NewBuffer()
	)

	err := encode.MarshalPackBody(bodyBuf, p.PackBody)

	if err != nil {
		return nil, err
	}

	p.PackHeader.Property.BodyByteLength = uint16(bodyBuf.Len())

	packHeaderData, err := p.PackHeader.MarshalBinary()

	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(packHeaderData); err != nil {
		return nil, err
	}

	if _, err := buf.Write(bodyBuf.Bytes()); err != nil {
		return nil, err
	}

	checksum, err := bytes.CalcChecksum(buf.Bytes())

	if err != nil {
		return nil, err
	}

	if err := buf.WriteByte(checksum); err != nil {
		return nil, err
	}

	escapedBytes, err := bytes.Escape(buf.Bytes())

	if err != nil {
		return nil, err
	}

	finalBuf.WriteByte(internal.IdentifyByte)

	if _, err := finalBuf.Write(escapedBytes); err != nil {
		return nil, err
	}

	finalBuf.WriteByte(internal.IdentifyByte)

	return finalBuf.Bytes(), nil
}

func (p *MessagePack[T]) UnmarshalBinary(buf []byte) error {
	var (
		reader         *bytes.Reader
		packBodyReader *bytes.Reader
		checksumGot    byte
	)

	if buf[0] != internal.IdentifyByte {
		return fmt.Errorf("invalid prefix byte 0x%.2X", buf[0])
	}

	if buf[len(buf)-1] != internal.IdentifyByte {
		return fmt.Errorf("invalid suffix byte 0x%.2X", buf[len(buf)-1])
	}

	buf = buf[1 : len(buf)-1]

	buf, err := bytes.Unescape(buf)

	if err != nil {
		return err
	}

	checksumGot, err = bytes.CalcChecksum(buf[0 : len(buf)-1])

	if err != nil {
		return err
	}

	reader = bytes.NewReader(buf)

	// read header, ( 12 or 12 + 4 bytes depends on is multiple package message)
	packHeaderData, err := reader.ReadBytes(16)

	if err != nil {
		return err
	}

	if err := p.PackHeader.UnmarshalBinary(packHeaderData); err != nil {
		return err
	}

	// is not a multiple package, reverse reader 4 bytes back because there's no package bytes
	if !p.PackHeader.Property.IsMultiplePackage {
		for i := 0; i < 4; i++ {
			_ = reader.UnreadByte()
		}
	}

	// update PackBody field from readed bytes to struct
	packBody, err := p.NewPackBodyFromMesgId()

	if err != nil {
		return err
	}

	// read bytes according header body data length
	packBodyData, err := reader.ReadBytes(int(p.PackHeader.Property.BodyByteLength))

	if err != nil {
		return err
	}

	packBodyReader = bytes.NewReader(packBodyData)

	if err := decode.UnmarshalPackBody(packBodyReader, packBody); err != nil {
		return err
	}

	// update checksum in message pack
	checksumWant, err := reader.ReadByte()

	if err != nil {
		return err
	}

	typedPackBody, ok := packBody.(T)

	if !ok {
		return fmt.Errorf("can't convert mesgBody %T as %T", packBody, p.PackBody)
	}

	p.Checksum = checksumWant
	p.PackBody = typedPackBody
	p.ChecksumValid = checksumWant == checksumGot

	return nil
}
