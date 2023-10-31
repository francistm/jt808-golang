package message

type Body0102_13 struct {
	AuthCode string
}

func (body *Body0102_13) MarshalBinary() ([]byte, error) {
	return []byte(body.AuthCode), nil
}

func (body *Body0102_13) UnmarshalBinary(data []byte) error {
	body.AuthCode = string(data)

	return nil
}
