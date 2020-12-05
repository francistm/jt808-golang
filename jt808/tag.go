package jt808

import (
	"strconv"
	"strings"
)

const (
	tagName = "jt808"

	tagEncodingAuto = "auto"
	tagEncodingBCD  = "bcd"
	tagEncodingNone = "none"
)

type parsedTag struct {
	fieldDataLength   int
	fieldDataEncoding string
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
		i, err := strconv.ParseUint(tags[0], 10, 8)

		if err != nil {
			return nil, err
		}

		t.fieldDataLength = int(i)

	case 2:
		i, err := strconv.ParseUint(tags[0], 10, 8)

		if err != nil {
			return nil, err
		}

		t.fieldDataLength = int(i)
		t.fieldDataEncoding = tags[1]
	}

	return t, nil
}
