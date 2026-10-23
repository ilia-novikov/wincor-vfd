package display

import "fmt"

var (
	clearDisplayCmd = []uint8{0x1B, 0x5B, 0x32, 0x4A}
)

func setCharsetCmd(charset Charset) []uint8 {
	return []uint8{0x1B, 0x52, (uint8)(charset)}
}

func setLineCmd(line int) ([]uint8, error) {
	if line != 0 && line != 1 {
		return nil, fmt.Errorf("only 2 lines are supported, got %d", line)
	}

	return []uint8{0x1B, 0x5B, uint8(0x31 + line), 0x3B, 0x31, 0x48}, nil
}
