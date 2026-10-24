package display

import "golang.org/x/text/encoding/charmap"

type Charset uint8

const (
	CharsetUS Charset = 0x00
	CharsetRU Charset = 0x35
)

func encodeText(charset Charset, text string) ([]uint8, error) {
	switch charset {
	case CharsetUS:
		return []uint8(text), nil

	case CharsetRU:
		encoded, err := charmap.CodePage866.NewEncoder().Bytes([]uint8(text))
		if err != nil {
			return nil, err
		}
		return encoded, nil

	default:
		return nil, &CharsetError{Charset: charset}
	}
}
