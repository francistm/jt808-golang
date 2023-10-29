package jt808

import (
	"strconv"
	"strings"
)

const (
	tagName = "jt808"

	tagEncodingAuto = "auto"
	tagEncodingBCD  = "bcd"
	tagEncodingGBK  = "gbk"
	tagEncodingNone = "raw"
)

type mesgTag struct {
	dataLength   int
	dataEncoding string
}

func (t *mesgTag) parseDataLength(s string) error {
	if len(s) == 0 {
		return nil
	}

	i, err := strconv.ParseUint(s, 10, 8)

	if err == nil {
		t.dataLength = int(i)
	}

	return err
}

func parseMesgTag(tag string) (*mesgTag, error) {
	t := &mesgTag{
		dataLength:   -1,
		dataEncoding: tagEncodingAuto,
	}

	if tag == "" {
		return t, nil
	}

	commaIndex := strings.IndexByte(tag, ',')

	if err := t.parseDataLength(tag[:commaIndex]); err != nil {
		return nil, err
	}

	if commaIndex > -1 {
		t.dataEncoding = tag[commaIndex+1:]
	}

	return t, nil
}
