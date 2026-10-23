package display

import "fmt"

var (
	clearDisplayCmd = []byte{0x1B, 0x5B, 0x32, 0x4A}
)

func setCharsetCmd(charset Charset) ([]byte, error) {
	// TODO: Move closer to Charset?
	switch charset {
	case US:
		return []byte{0x1B, 0x52, 0x00}, nil
	case RU:
		return []byte{0x1B, 0x52, 0x35}, nil
	default:
		return nil, &CharsetError{Charset: charset}
	}
}

func setLineCmd(line int) ([]byte, error) {
	if line != 0 && line != 1 {
		return nil, fmt.Errorf("only 2 lines are supported, got %d", line)
	}

	return []byte{0x1B, 0x5B, byte(0x31 + line), 0x3B, 0x31, 0x48}, nil
}
