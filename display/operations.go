package display

// PrintText encodes and writes text at the current cursor position.
func (d *Display) PrintText(text string) error {
	msg, err := encodeText(d.charset, text)
	if err != nil {
		return err
	}

	return d.WriteData(msg)
}

// SendBackspace moves the cursor one space to the left. If there is a character
// in the position to which the cursor moves, it is not deleted.
func (d *Display) SendBackspace() error {
	return d.WriteData(backspaceCmd)
}

// SendLineFeed moves the cursor down a line, the column position remains the same.
// If the cursors is already on the last line, the contents of the last line is
// copied to the first line and the last line is deleted.
func (d *Display) SendLineFeed() error {
	return d.WriteData(lineFeedCmd)
}

// SendCarriageReturn moves the cursor to the start of the current line.
func (d *Display) SendCarriageReturn() error {
	return d.WriteData(carriageReturnCmd)
}

// ClearDisplay clears the whole display. The cursor position remains unchanged.
func (d *Display) ClearDisplay() error {
	return d.WriteData(clearDisplayCmd)
}

// SetCharset selects the display's configured charset.
func (d *Display) SetCharset() error {
	cmd := setCharsetCmd(d.charset)
	return d.WriteData(cmd)
}

// MoveCursor moves the cursor to the given 0-based column/row.
func (d *Display) MoveCursor(x, y int) error {
	return d.WriteData(moveCursorCmd(x, y))
}

// DeleteToEndOfLine erases text to the end the current line, cursor position
// included. The cursor position remains unchanged.
func (d *Display) DeleteToEndOfLine() error {
	return d.WriteData(deleteToEndOfLineCmd)
}
