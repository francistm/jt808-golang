package tag

import (
	"strconv"
	"strings"
)

const (
	Name = "jt808"

	EncodingAuto = "auto"
	EncodingBCD  = "bcd"
	EncodingGBK  = "gbk"
	EncodingRaw  = "raw"
)

type MesgTag struct {
	Length   int
	Encoding string
}

func NewMesgTag(s string) (*MesgTag, error) {
	t := &MesgTag{
		Length:   0,
		Encoding: EncodingAuto,
	}

	if s == "" {
		return t, nil
	}

	var (
		length   string
		encoding string
	)

	switch strings.Count(s, ",") {
	case 0:
		length = s

	case 1:
		i := strings.IndexByte(s, ',')

		length = s[:i]
		encoding = s[i+1:]
	}

	if length != "" {
		lenInt, err := strconv.ParseInt(length, 10, 8)

		if err != nil {
			return nil, err
		}

		t.Length = int(lenInt)
	}

	if encoding != "" {
		t.Encoding = encoding
	}

	return t, nil
}
