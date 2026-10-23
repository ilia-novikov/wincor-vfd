package display

var (
	clearDisplayCmd = []byte{0x1B, 0x5B, 0x32, 0x4A}
)

func setCharsetCmd(charset Charset) ([]byte, error) {
	// TODO: Move closer to Charset?
	switch charset {
	case US:
		return []byte{0x1b, 0x52, 0x00}, nil
	case RU:
		return []byte{0x1b, 0x52, 0x35}, nil
	default:
		return nil, &CharsetError{Charset: charset}
	}
}
