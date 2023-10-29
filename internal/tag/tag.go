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
		Length:   -1,
		Encoding: EncodingAuto,
	}

	if s == "" {
		return t, nil
	}

	switch strings.Count(s, ",") {
	case 0:
		lenInt, err := strconv.ParseUint(s, 10, 8)

		if err != nil {
			return nil, err
		}

		if lenInt > 0 {
			t.Length = int(lenInt)
		}
	case 1:
		i := strings.IndexByte(s, ',')
		lenString := s[:i]
		enc := s[i+1:]

		if lenString != "" {
			lenInt, err := strconv.ParseUint(lenString, 10, 8)

			if err != nil {
				return nil, err
			}

			if lenInt > 0 {
				t.Length = int(lenInt)
			}
		}

		if enc != "" {
			t.Encoding = enc
		}
	}

	return t, nil
}
