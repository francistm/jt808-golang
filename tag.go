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
	tagEncodingNone = "none"
)

type parsedTag struct {
	fieldDataLength   int
	fieldDataEncoding string
}

func (ptr *parsedTag) parseFieldDataLength(s string) error {
	if len(s) == 0 {
		return nil
	}

	i, err := strconv.ParseUint(s, 10, 8)

	if err == nil {
		ptr.fieldDataLength = int(i)
	}

	return err
}

func parseTag(tag string) (*parsedTag, error) {
	t := &parsedTag{
		fieldDataLength:   -1,
		fieldDataEncoding: tagEncodingAuto,
	}

	if tag == "" {
		return t, nil
	}

	tags := strings.Split(tag, ",")

	switch len(tags) {
	case 1:
		if err := t.parseFieldDataLength(tags[0]); err != nil {
			return nil, err
		}

	case 2:
		if err := t.parseFieldDataLength(tags[0]); err != nil {
			return nil, err
		}

		t.fieldDataEncoding = tags[1]
	}

	return t, nil
}
