// Code generated by generator/decoder, DO NOT MODIFY MANUALLY

package message

import "fmt"

func (*Body0001) MesgId() uint16 {
	return uint16(0x1)
}

func (*Body0002) MesgId() uint16 {
	return uint16(0x2)
}

func (*Body0003) MesgId() uint16 {
	return uint16(0x3)
}

func (*Body0004) MesgId() uint16 {
	return uint16(0x4)
}

func (*Body0005) MesgId() uint16 {
	return uint16(0x5)
}

func (*Body0100_13) MesgId() uint16 {
	return uint16(0x100)
}

func (*Body0100_19) MesgId() uint16 {
	return uint16(0x100)
}

func (*Body0102_13) MesgId() uint16 {
	return uint16(0x102)
}

func (*Body0102_19) MesgId() uint16 {
	return uint16(0x102)
}

func (*Body0200) MesgId() uint16 {
	return uint16(0x200)
}

func (*Body0801) MesgId() uint16 {
	return uint16(0x801)
}

func (*Body8001) MesgId() uint16 {
	return uint16(0x8001)
}

func (*Body8003) MesgId() uint16 {
	return uint16(0x8003)
}

func (*Body8004) MesgId() uint16 {
	return uint16(0x8004)
}

func (*Body8100) MesgId() uint16 {
	return uint16(0x8100)
}

func (m *MessagePack[T]) NewPackBodyFromMesgId() (MesgBody, error) {
	if m.PackHeader.Package != nil {
		return new(PartialPackBody), nil
	} else {
		switch m.PackHeader.MessageID {
		case uint16(0x1):
			return new(Body0001), nil
		case uint16(0x2):
			return new(Body0002), nil
		case uint16(0x3):
			return new(Body0003), nil
		case uint16(0x4):
			return new(Body0004), nil
		case uint16(0x5):
			return new(Body0005), nil
		case uint16(0x100):
			if m.PackHeader.Property.Version == Version2013 {
				return new(Body0100_13), nil
			} else if m.PackHeader.Property.Version == Version2019 {
				return new(Body0100_19), nil
			} else {
				return nil, fmt.Errorf("unsupport protocol version: %d", m.PackHeader.Property.Version)
			}
		case uint16(0x102):
			if m.PackHeader.Property.Version == Version2013 {
				return new(Body0102_13), nil
			} else if m.PackHeader.Property.Version == Version2019 {
				return new(Body0102_19), nil
			} else {
				return nil, fmt.Errorf("unsupport protocol version: %d", m.PackHeader.Property.Version)
			}
		case uint16(0x200):
			return new(Body0200), nil
		case uint16(0x801):
			return new(Body0801), nil
		case uint16(0x8001):
			return new(Body8001), nil
		case uint16(0x8003):
			return new(Body8003), nil
		case uint16(0x8004):
			return new(Body8004), nil
		case uint16(0x8100):
			return new(Body8100), nil
		default:
			return nil, fmt.Errorf("unsupported messageId: 0x%.4X", m.PackHeader.MessageID)
		}
	}
}
