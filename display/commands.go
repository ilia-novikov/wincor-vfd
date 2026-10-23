package display

var (
	backspaceCmd      = []uint8{0x08}
	lineFeedCmd       = []uint8{0x0A}
	carriageReturnCmd = []uint8{0x0D}

	clearDisplayCmd      = []uint8{0x1B, 0x5B, 0x32, 0x4A}
	deleteToEndOfLineCmd = []uint8{0x1B, 0x5B, 0x30, 0x4B}
)

func setCharsetCmd(charset Charset) []uint8 {
	return []uint8{0x1B, 0x52, (uint8)(charset)}
}

func moveCursorCmd(x, y int) []uint8 {
	return []uint8{0x1B, 0x5B, uint8(0x31 + y), 0x3B, uint8(0x31 + x), 0x48}
}
