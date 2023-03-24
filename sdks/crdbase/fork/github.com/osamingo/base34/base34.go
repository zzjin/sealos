package base34

import (
	"encoding/binary"
	"errors"
)

const (
	// A CharacterLength defines source characters length.
	CharacterLength = 34
	// A StandardSource defines standard characters of base34.
	StandardSource = "123456789abcdefghijkmnopqrstuvwxyz"
)

var (
	// ErrInvalidArgument expresses invalid argument error.
	ErrInvalidArgument = errors.New("base34: invalid argument")
	// ErrUnknownCharacter expresses contains an unknown character error.
	ErrUnknownCharacter = errors.New("base34: contains an unknown character")
)

// An Encoder implements encoding and decoding of base34.
type Encoder struct {
	encode    [CharacterLength]byte
	decodeMap [256]int
}

// MustNewEncoder returns new base34.Encoder.
func MustNewEncoder(source string) *Encoder {
	enc, err := NewEncoder(source)
	if err != nil {
		panic(err)
	}

	return enc
}

// NewEncoder returns new base34.Encoder.
func NewEncoder(source string) (*Encoder, error) {
	if len(source) != CharacterLength {
		return nil, ErrInvalidArgument
	}

	enc := new(Encoder)
	for i := range enc.decodeMap {
		enc.decodeMap[i] = -1
	}

	for i := range source {
		enc.encode[i] = source[i]
		enc.decodeMap[enc.encode[i]] = i
	}

	return enc, nil
}

// Encode returns encoded string by base34.
func (enc *Encoder) Encode(id uint64) string {
	if id == 0 {
		return string(enc.encode[:1])
	}

	bin := make([]byte, 0, binary.MaxVarintLen64)
	for id > 0 {
		bin = append(bin, enc.encode[id%CharacterLength])
		id /= CharacterLength
	}

	for i, j := 0, len(bin)-1; i < j; i, j = i+1, j-1 {
		bin[i], bin[j] = bin[j], bin[i]
	}

	return string(bin)
}

// Decode returns decoded unsigned int64 by base34.
func (enc *Encoder) Decode(id string) (uint64, error) {
	if id == "" {
		return 0, ErrInvalidArgument
	}
	n := uint64(0)

	for _, v := range id {
		u := enc.decodeMap[v]
		if u < 0 {
			return 0, ErrUnknownCharacter
		}

		n = n*CharacterLength + uint64(u)
	}

	return n, nil
}
